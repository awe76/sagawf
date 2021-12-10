package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/awe76/sagawf/workflow"
	"go-micro.dev/v4/broker"

	pb "github.com/awe76/sagawf/proto"
)

type Sagawf struct {
	cache    workflow.Cache
	producer workflow.Producer
	handler  map[int]chan workflow.WorkflowPayload
}

func getWorkflowKey(id int) string {
	return fmt.Sprintf("workflow:definition:%d", id)
}

func NewSagawf() (*Sagawf, error) {
	cache := workflow.NewCache()
	producer := workflow.NewProducer()
	err := producer.Init()
	if err != nil {
		return nil, err
	}

	err = producer.Connect()
	if err != nil {
		return nil, err
	}

	handler := make(map[int]chan workflow.WorkflowPayload)

	result := Sagawf{
		cache:    cache,
		producer: producer,
		handler:  handler,
	}

	_, err = broker.Subscribe(workflow.WORKFLOW_OPERATION_START, func(p broker.Event) error {
		var op workflow.OperationPayload
		err := json.Unmarshal(p.Message().Body, &op)
		if err != nil {
			return err
		}

		if op.IsRollback {
			fmt.Printf("%s operation rollback is started\n", op.Operation.Name)
		} else {
			fmt.Printf("%s operation is started\n", op.Operation.Name)
		}

		rand.Seed(time.Now().UnixNano())

		pause := rand.Intn(100)
		// sleep for some random time
		time.Sleep(time.Duration(pause) * time.Millisecond)

		op.Payload = rand.Float32()

		// randomly complete or fault the operation
		if op.IsRollback || rand.Float32() < 0.8 {
			return producer.SendMessage(workflow.WORKFLOW_OPERATION_COMPLETED, op)
		} else {
			return producer.SendMessage(workflow.WORKFLOW_OPERATION_FAILED, op)
		}
	})

	if err != nil {
		return nil, err
	}

	_, err = broker.Subscribe(workflow.WORKFLOW_OPERATION_COMPLETED, func(p broker.Event) error {
		var op workflow.OperationPayload
		err := json.Unmarshal(p.Message().Body, &op)
		if err != nil {
			return err
		}

		if op.IsRollback {
			fmt.Printf("%s operation rollback is completed\n", op.Operation.Name)
		} else {
			fmt.Printf("%s operation is completed\n", op.Operation.Name)
		}

		proc := result.CreateProcessor()
		w, err := result.GetWorkflow(op.ID)
		if err != nil {
			return err
		}

		return proc.OnComplete(w, op)
		return nil
	})

	if err != nil {
		return nil, err
	}

	_, err = broker.Subscribe(workflow.WORKFLOW_OPERATION_FAILED, func(p broker.Event) error {
		var op workflow.OperationPayload
		err := json.Unmarshal(p.Message().Body, &op)
		if err != nil {
			return err
		}

		fmt.Printf("%s operation is failed\n", op.Operation.Name)
		proc := result.CreateProcessor()
		w, err := result.GetWorkflow(op.ID)
		if err != nil {
			return err
		}
		return proc.OnFailure(w, op)
	})

	if err != nil {
		return nil, err
	}

	_, err = broker.Subscribe(workflow.WORKFLOW_COMPLETED, func(p broker.Event) error {
		var w workflow.WorkflowPayload
		err := json.Unmarshal(p.Message().Body, &w)
		if err != nil {
			return err
		}

		fmt.Printf("%s %d workflow is completed\n", w.Name, w.ID)
		fmt.Printf("workflow state: %v\n", w.Data)

		if ch, found := handler[w.ID]; found {
			ch <- w
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	_, err = broker.Subscribe(workflow.WORKFLOW_ROLLBACKED, func(p broker.Event) error {
		var w workflow.WorkflowPayload
		err := json.Unmarshal(p.Message().Body, &w)
		if err != nil {
			return err
		}

		fmt.Printf("%s %d workflow is rollbacked\n", w.Name, w.ID)
		fmt.Printf("workflow state: %v\n", w.Data)

		if ch, found := handler[w.ID]; found {
			ch <- w
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (e *Sagawf) CreateProcessor() workflow.Processor {
	return workflow.NewProcessor(e.cache, e.producer)
}

func (e *Sagawf) ReserveID() (int, error) {
	id, err := workflow.ReserveID("workflow:index", e.cache)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (e *Sagawf) GetWorkflow(id int) (workflow.Workflow, error) {
	ctx := context.Background()
	key := getWorkflowKey(id)

	var w workflow.Workflow
	raw, err := e.cache.Get(ctx, key)
	if err != nil {
		return w, err
	}

	json.Unmarshal([]byte(raw), &w)
	return w, nil
}

func (e *Sagawf) RegisterHandler(id int) chan workflow.WorkflowPayload {
	result := make(chan workflow.WorkflowPayload)
	e.handler[id] = result

	return result
}

func (e *Sagawf) SetWorkflow(id int, w workflow.Workflow) error {
	ctx := context.Background()
	key := getWorkflowKey(id)
	return e.cache.Set(ctx, key, w)
}

func (e *Sagawf) RunWorkflow(ctx context.Context, req *pb.WorkflowRequest, rsp *pb.WorkflowResponse) error {
	proc := e.CreateProcessor()
	id, err := e.ReserveID()

	if err != nil {
		return err
	}

	targetCh := e.RegisterHandler(id)

	w := workflow.ToWorkflow(req)
	err = e.SetWorkflow(id, w)
	if err != nil {
		return err
	}

	proc.StartWorkflow(w, id)
	if err != nil {
		return err
	}

	response := <-targetCh

	rsp.WorkflowRef = &pb.WorkflowRef{
		Id:         int64(id),
		Name:       req.Name,
		IsRollback: response.IsRollback,
	}

	rsp.State = make(map[string]*pb.State)

	for s, v := range response.Data {
		st := pb.State{
			State: make(map[string]string),
		}
		rsp.State[s] = &st
		for op, state := range v {
			opValue, err := json.Marshal(state)
			if err != nil {
				return err
			}

			st.State[op] = string(opValue)
		}
	}
	return nil
}
