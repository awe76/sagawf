package workflow

type OperationPayload struct {
	ID         int
	IsRollback bool
	Name       string
	Operation  Operation
	Payload    interface{}
}
