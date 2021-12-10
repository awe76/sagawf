package workflow

import (
	"context"
	"encoding/json"
	"fmt"
)

type state struct {
	ID         int
	IsRollback bool
	Completed  bool
	Done       map[string]Operation
	InProgress map[string]Operation
	Data       map[string]map[string]interface{}
}

func (s *state) getCacheKey() string {
	return fmt.Sprintf("workflow:state:%v", s.ID)
}

func (s *state) init(cache Cache, start string, payload interface{}) error {
	ctx := context.Background()

	s.IsRollback = false
	s.Completed = false
	s.Done = make(map[string]Operation)
	s.InProgress = make(map[string]Operation)
	s.Data = make(map[string]map[string]interface{})
	s.setData(start, "input", payload)

	key := s.getCacheKey()
	cache.Set(ctx, key, s)
	return nil
}

func (s *state) setData(vertex string, operation string, payload interface{}) {
	ops, found := s.Data[vertex]
	if !found {
		ops = make(map[string]interface{})
		s.Data[vertex] = ops
	}

	ops[operation] = payload
}

func (s *state) update(cache Cache, update func(*state)) error {
	ctx := context.Background()

	key := s.getCacheKey()
	rawState, err := cache.Get(ctx, key)

	if err != nil {
		return err
	}

	json.Unmarshal([]byte(rawState), s)
	if err != nil {
		return err
	}

	update(s)
	err = cache.Set(ctx, key, s)
	return err
}
