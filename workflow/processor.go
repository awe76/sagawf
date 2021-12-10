package workflow

const (
	WORKFLOW_OPERATION_START     = "wfos"
	WORKFLOW_OPERATION_COMPLETED = "wfoc"
	WORKFLOW_OPERATION_FAILED    = "wfof"
	WORKFLOW_COMPLETED           = "wfc"
	WORKFLOW_ROLLBACKED          = "wfr"
	WORKFLOW_START               = "wfs"
)

type RouteMap struct {
	Route map[string][]Operation
}

type processor struct {
	cache    Cache
	producer Producer
	workflow Workflow
	state    state
}

func NewProcessor(cache Cache, producer Producer) Processor {
	return &processor{
		cache:    cache,
		producer: producer,
	}
}

type Processor interface {
	StartWorkflow(w Workflow, id int) error
	OnComplete(w Workflow, op OperationPayload) error
	OnFailure(w Workflow, op OperationPayload) error
}

func (p *processor) StartWorkflow(w Workflow, id int) error {
	p.workflow = w
	p.state = state{
		ID: id,
	}
	err := p.state.init(p.cache, w.Start, w.Payload)
	if err != nil {
		return err
	}

	t := createDirectTracer(w, p.state, p.endWorkflow, p.spawnOperation)
	return t.resolveWorkflow(w.Start)
}

func (p *processor) OnComplete(w Workflow, op OperationPayload) error {
	p.workflow = w

	p.state = state{
		ID: op.ID,
	}
	err := p.state.update(p.cache, func(s *state) {
		removeOp(s.InProgress, op.Operation, op.IsRollback)
		addOp(s.Done, op.Operation, op.IsRollback)
		s.setData(op.Operation.To, op.Operation.Name, op.Payload)
	})

	if err != nil {
		return err
	}

	if p.state.IsRollback {
		t := createReverseTracer(w, p.state, p.endWorkflow, p.spawnOperation)
		return t.resolveWorkflow(w.End)
	} else {
		t := createDirectTracer(w, p.state, p.endWorkflow, p.spawnOperation)
		return t.resolveWorkflow(w.Start)
	}
}

func (p *processor) OnFailure(w Workflow, op OperationPayload) error {
	p.workflow = w

	p.state = state{
		ID: op.ID,
	}
	err := p.state.update(p.cache, func(s *state) {
		removeOp(s.InProgress, op.Operation, false)

		s.IsRollback = true
	})

	if err != nil {
		return err
	}

	t := createReverseTracer(w, p.state, p.endWorkflow, p.spawnOperation)
	return t.resolveWorkflow(w.End)
}

func (p *processor) spawnOperation(op Operation) error {

	data := p.state.Data[op.From]

	payload := OperationPayload{
		ID:         p.state.ID,
		IsRollback: p.state.IsRollback,
		Name:       p.workflow.Name,
		Operation:  op,
		Payload:    data,
	}

	err := p.state.update(p.cache, func(s *state) {
		addOp(s.InProgress, op, p.state.IsRollback)
	})

	if err != nil {
		return err
	}

	return p.producer.SendMessage(WORKFLOW_OPERATION_START, payload)
}

func (p *processor) endWorkflow() error {
	if !p.state.Completed {
		err := p.state.update(p.cache, func(s *state) {
			s.Completed = true
		})
		if err != nil {
			return err
		}

		payload := WorkflowPayload{
			ID:         p.state.ID,
			IsRollback: p.state.IsRollback,
			Name:       p.workflow.Name,
			Data:       p.state.Data,
		}

		topic := WORKFLOW_COMPLETED
		if p.state.IsRollback {
			topic = WORKFLOW_ROLLBACKED
		}

		return p.producer.SendMessage(topic, payload)
	}

	return nil
}
