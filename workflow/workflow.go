package workflow

import (
	pb "github.com/awe76/sagawf/proto"
)

type Workflow struct {
	Name       string      `json:"name"`
	Start      string      `json:"start"`
	End        string      `json:"end"`
	Operations []Operation `json:"operations"`
	Payload    interface{} `json:"payload"`
}

func (w *Workflow) toPayload(id int, isReversion bool, data map[string]map[string]interface{}) WorkflowPayload {
	return WorkflowPayload{
		ID:         id,
		Name:       w.Name,
		Data:       data,
		IsRollback: isReversion,
	}
}

func ToWorkflow(req *pb.WorkflowRequest) Workflow {
	return Workflow{
		Name:       req.Name,
		Start:      req.Start,
		End:        req.End,
		Operations: toOperations(req.Operations),
		Payload:    req.Payload,
	}
}

func toOperations(ops []*pb.Operation) []Operation {
	result := []Operation{}
	for _, op := range ops {
		result = append(result, Operation{
			Name: op.Name,
			From: op.From,
			To:   op.To,
		})
	}

	return result
}
