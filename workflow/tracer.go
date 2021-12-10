package workflow

type route = map[string][]Operation
type getOperationKey = func(op Operation) string

func getFrom(op Operation) string {
	return op.From
}

func getTo(op Operation) string {
	return op.To
}

func hasOp(m map[string]Operation, op Operation, isRollback bool) bool {
	key := op.getKey(isRollback)
	_, found := m[key]
	return found
}

func addOp(m map[string]Operation, op Operation, isRollback bool) {
	key := op.getKey(isRollback)
	m[key] = op
}

func removeOp(m map[string]Operation, op Operation, isRollback bool) {
	key := op.getKey(isRollback)
	delete(m, key)
}

func allMatched(current string, relation route, isMatched func(op Operation) bool) bool {
	// get all related operations
	if ops, found := relation[current]; found {
		// for each operation
		for _, op := range ops {
			if !isMatched(op) {
				return false
			}
		}
	}

	return true
}

func createRoute(operations []Operation, getKey getOperationKey) route {
	result := make(map[string][]Operation)
	for _, op := range operations {
		key := getKey(op)
		if next, found := result[key]; found {
			result[key] = append(next, op)
		} else {
			result[key] = append([]Operation{}, op)
		}
	}

	return result
}

type tracer struct {
	isReady        func(current string) bool
	isFinished     func(current string) bool
	getNext        func(current string) (ops []Operation, found bool)
	isProcessed    func(op Operation) bool
	canBeSpawned   func(op Operation) bool
	getNextVertex  func(op Operation) string
	endWorkflow    func() error
	spawnOperation func(op Operation) error
}

func createDirectTracer(
	w Workflow,
	s state,
	endWorkflow func() error,
	spawnOperation func(op Operation) error,
) *tracer {
	from := createRoute(w.Operations, getFrom)
	to := createRoute(w.Operations, getTo)

	isMatched := func(op Operation) bool {
		return hasOp(s.Done, op, false)
	}

	return &tracer{
		isReady: func(current string) bool {
			return allMatched(current, to, isMatched)
		},
		isFinished: func(current string) bool {
			return current == w.End
		},
		getNext: func(current string) ([]Operation, bool) {
			ops, found := from[current]
			return ops, found
		},
		isProcessed: func(op Operation) bool {
			return hasOp(s.Done, op, false)
		},
		canBeSpawned: func(op Operation) bool {
			return !hasOp(s.InProgress, op, false)
		},
		getNextVertex: func(op Operation) string {
			return op.To
		},
		endWorkflow: endWorkflow,
		spawnOperation: func(op Operation) error {
			return spawnOperation(op)
		},
	}
}

func createReverseTracer(
	w Workflow,
	s state,
	endWorkflow func() error,
	spawnOperation func(op Operation) error,
) *tracer {
	from := createRoute(w.Operations, getFrom)
	to := createRoute(w.Operations, getTo)

	isMatched := func(op Operation) bool {
		done := hasOp(s.Done, op, false)
		return !hasOp(s.InProgress, op, false) && !hasOp(s.InProgress, op, true) && (!done || (done && hasOp(s.Done, op, true)))
	}

	return &tracer{
		isReady: func(current string) bool {
			return allMatched(current, from, isMatched)
		},
		isFinished: func(current string) bool {
			return current == w.Start
		},
		getNext: func(current string) ([]Operation, bool) {
			ops, found := to[current]
			return ops, found
		},
		isProcessed: func(op Operation) bool {
			done := hasOp(s.Done, op, false)
			return !done || (done && hasOp(s.Done, op, true))
		},
		canBeSpawned: func(op Operation) bool {
			return !hasOp(s.InProgress, op, true)
		},
		getNextVertex: func(op Operation) string {
			return op.From
		},
		endWorkflow: endWorkflow,
		spawnOperation: func(op Operation) error {
			return spawnOperation(op)
		},
	}
}

func (t *tracer) resolveWorkflow(current string) error {
	// current vertex is ready for resolution
	if t.isReady(current) {
		if t.isFinished(current) {
			err := t.endWorkflow()
			if err != nil {
				return err
			}
		} else {
			// get next operations
			if ops, found := t.getNext(current); found {
				// for each next operation
				for _, op := range ops {
					if t.isProcessed(op) {
						// if operation has been already processed continue resolution of the next vertex
						nextVertex := t.getNextVertex(op)
						err := t.resolveWorkflow(nextVertex)
						if err != nil {
							return err
						}
					} else if t.canBeSpawned(op) {
						// if operation can be spawned do it
						err := t.spawnOperation(op)
						if err != nil {
							return err
						}
					}
				}
			}
		}
	}

	return nil
}
