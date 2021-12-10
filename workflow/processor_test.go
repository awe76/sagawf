package workflow

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessor(t *testing.T) {

	ops := []Operation{
		{
			Name: "op1",
			From: "s1",
			To:   "s2",
		},
		{
			Name: "op2",
			From: "s1",
			To:   "s3",
		},
		{
			Name: "op3",
			From: "s3",
			To:   "s2",
		},
	}

	defaultWorkflow := Workflow{
		Name:       "default workflow",
		Start:      "s1",
		End:        "s2",
		Operations: ops,
	}

	type step struct {
		action   func(t *testing.T, w Workflow, p *processor)
		validate func(t *testing.T, w Workflow, p *ProducerMock)
	}

	var tests = map[string]struct {
		w     Workflow
		steps []step
	}{
		"default workflow is completed": {
			w: defaultWorkflow,
			steps: []step{
				{
					action: func(t *testing.T, w Workflow, p *processor) {
						err := p.StartWorkflow(w, 1)
						assert.NoError(t, err)
					},
					validate: func(t *testing.T, w Workflow, p *ProducerMock) {
						payload := make(map[string]interface{})
						payload["input"] = nil
						op1 := ops[0].toPayload(1, w, false, payload)
						assert.True(t, p.Has(WORKFLOW_OPERATION_START, op1))

						op2 := ops[1].toPayload(1, w, false, payload)
						assert.True(t, p.Has(WORKFLOW_OPERATION_START, op2))
					},
				},
				{
					action: func(t *testing.T, w Workflow, p *processor) {
						op1 := ops[0].toPayload(1, w, false, nil)
						assert.NoError(t, p.OnComplete(w, op1))

						op2 := ops[1].toPayload(1, w, false, nil)
						assert.NoError(t, p.OnComplete(w, op2))
					},
					validate: func(t *testing.T, w Workflow, p *ProducerMock) {
						payload := make(map[string]interface{})
						payload["op2"] = nil
						op3 := ops[2].toPayload(1, w, false, payload)
						assert.True(t, p.Has(WORKFLOW_OPERATION_START, op3))
					},
				},
				{
					action: func(t *testing.T, w Workflow, p *processor) {
						op3 := ops[2].toPayload(1, w, false, nil)
						assert.NoError(t, p.OnComplete(w, op3))
					},
					validate: func(t *testing.T, w Workflow, p *ProducerMock) {
						data := make(map[string]map[string]interface{})
						data["s1"] = make(map[string]interface{})
						data["s2"] = make(map[string]interface{})
						data["s3"] = make(map[string]interface{})

						data["s1"]["input"] = nil
						data["s2"]["op1"] = nil
						data["s2"]["op3"] = nil
						data["s3"]["op2"] = nil
						wp := w.toPayload(1, false, data)
						assert.True(t, p.Has(WORKFLOW_COMPLETED, wp))
					},
				},
			},
		},
		"default workflow is rollbacked": {
			w: defaultWorkflow,
			steps: []step{
				{
					action: func(t *testing.T, w Workflow, p *processor) {
						err := p.StartWorkflow(w, 1)
						assert.NoError(t, err)
					},
					validate: func(t *testing.T, w Workflow, p *ProducerMock) {
						payload := make(map[string]interface{})
						payload["input"] = nil
						op1 := ops[0].toPayload(1, w, false, payload)
						assert.True(t, p.Has(WORKFLOW_OPERATION_START, op1))

						op2 := ops[1].toPayload(1, w, false, payload)
						assert.True(t, p.Has(WORKFLOW_OPERATION_START, op2))
					},
				},
				{
					action: func(t *testing.T, w Workflow, p *processor) {
						op1 := ops[0].toPayload(1, w, false, nil)
						assert.NoError(t, p.OnComplete(w, op1))

						op2 := ops[1].toPayload(1, w, false, nil)
						assert.NoError(t, p.OnComplete(w, op2))
					},
					validate: func(t *testing.T, w Workflow, p *ProducerMock) {
						payload := make(map[string]interface{})
						payload["op2"] = nil
						op3 := ops[2].toPayload(1, w, false, payload)
						assert.True(t, p.Has(WORKFLOW_OPERATION_START, op3))
					},
				},
				{
					action: func(t *testing.T, w Workflow, p *processor) {
						op3 := ops[2].toPayload(1, w, false, nil)
						assert.NoError(t, p.OnFailure(w, op3))
					},
					validate: func(t *testing.T, w Workflow, p *ProducerMock) {
						payload := make(map[string]interface{})
						payload["input"] = nil
						op1 := ops[0].toPayload(1, w, true, payload)
						assert.True(t, p.Has(WORKFLOW_OPERATION_START, op1))

						op2 := ops[1].toPayload(1, w, true, payload)
						assert.True(t, p.Has(WORKFLOW_OPERATION_START, op2))
					},
				},
				{
					action: func(t *testing.T, w Workflow, p *processor) {
						op1 := ops[0].toPayload(1, w, true, nil)
						assert.NoError(t, p.OnComplete(w, op1))

						op2 := ops[1].toPayload(1, w, true, nil)
						assert.NoError(t, p.OnComplete(w, op2))
					},
					validate: func(t *testing.T, w Workflow, p *ProducerMock) {
						data := make(map[string]map[string]interface{})
						data["s1"] = make(map[string]interface{})
						data["s2"] = make(map[string]interface{})
						data["s3"] = make(map[string]interface{})

						data["s1"]["input"] = nil
						data["s2"]["op1"] = nil
						data["s3"]["op2"] = nil
						wp := w.toPayload(1, true, data)
						assert.True(t, p.Has(WORKFLOW_ROLLBACKED, wp))
					},
				},
			},
		},
		"default workflow op2 is failed": {
			w: defaultWorkflow,
			steps: []step{
				{
					action: func(t *testing.T, w Workflow, p *processor) {
						err := p.StartWorkflow(w, 1)
						assert.NoError(t, err)
					},
					validate: func(t *testing.T, w Workflow, p *ProducerMock) {
						payload := make(map[string]interface{})
						payload["input"] = nil
						op1 := ops[0].toPayload(1, w, false, payload)
						assert.True(t, p.Has(WORKFLOW_OPERATION_START, op1))

						op2 := ops[1].toPayload(1, w, false, payload)
						assert.True(t, p.Has(WORKFLOW_OPERATION_START, op2))
					},
				},
				{
					action: func(t *testing.T, w Workflow, p *processor) {
						op1 := ops[0].toPayload(1, w, false, nil)
						assert.NoError(t, p.OnComplete(w, op1))

						op2 := ops[1].toPayload(1, w, false, nil)
						assert.NoError(t, p.OnFailure(w, op2))
					},
					validate: func(t *testing.T, w Workflow, p *ProducerMock) {
						payload := make(map[string]interface{})
						payload["input"] = nil
						op1 := ops[0].toPayload(1, w, true, payload)
						assert.True(t, p.Has(WORKFLOW_OPERATION_START, op1))
					},
				},
			},
		},
		"default workflow op1 is failed": {
			w: defaultWorkflow,
			steps: []step{
				{
					action: func(t *testing.T, w Workflow, p *processor) {
						err := p.StartWorkflow(w, 1)
						assert.NoError(t, err)
					},
					validate: func(t *testing.T, w Workflow, p *ProducerMock) {
						payload := make(map[string]interface{})
						payload["input"] = nil

						op1 := ops[0].toPayload(1, w, false, payload)
						assert.True(t, p.Has(WORKFLOW_OPERATION_START, op1))

						op2 := ops[1].toPayload(1, w, false, payload)
						assert.True(t, p.Has(WORKFLOW_OPERATION_START, op2))
					},
				},
				{
					action: func(t *testing.T, w Workflow, p *processor) {
						op1 := ops[0].toPayload(1, w, false, nil)
						assert.NoError(t, p.OnFailure(w, op1))

						op2 := ops[1].toPayload(1, w, false, nil)
						assert.NoError(t, p.OnComplete(w, op2))
					},
					validate: func(t *testing.T, w Workflow, p *ProducerMock) {
						payload := make(map[string]interface{})
						payload["input"] = nil
						op2 := ops[1].toPayload(1, w, true, payload)
						assert.True(t, p.Has(WORKFLOW_OPERATION_START, op2))
					},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			cache := NewCacheMock()
			producer := NewProducerMock()

			for _, step := range tc.steps {
				proc := &processor{
					cache:    cache,
					producer: producer,
				}
				step.action(t, tc.w, proc)
				step.validate(t, tc.w, producer)
			}
		})
	}

}
