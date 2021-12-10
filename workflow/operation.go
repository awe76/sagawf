package workflow

import "fmt"

type Operation struct {
	Name string `json:"name"`
	From string `json:"from"`
	To   string `json:"to"`
}

func (op *Operation) getKey(isRollback bool) string {
	return fmt.Sprintf("%s:%s:%s:%v", op.Name, op.From, op.To, isRollback)
}

func (op *Operation) toPayload(id int, w Workflow, isRollback bool, payload interface{}) OperationPayload {
	return OperationPayload{
		ID:         id,
		Name:       w.Name,
		IsRollback: isRollback,
		Payload:    payload,
		Operation:  *op,
	}
}
