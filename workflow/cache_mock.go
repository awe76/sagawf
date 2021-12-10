package workflow

import (
	"context"
	"encoding/json"

	mc "go-micro.dev/v4/cache"
)

type CacheMock struct {
	values map[string]string
}

func NewCacheMock() *CacheMock {
	return &CacheMock{
		values: make(map[string]string),
	}
}

func (c *CacheMock) Set(ctx context.Context, key string, value interface{}) error {
	rawValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	c.values[key] = string(rawValue)
	return nil
}

func (c *CacheMock) Get(ctx context.Context, key string) (string, error) {
	value, found := c.values[key]
	if !found {
		return "", mc.ErrKeyNotFound
	}

	return value, nil
}

func (c *CacheMock) Remove(ctx context.Context, key string) error {
	delete(c.values, key)

	return nil
}

func (c *CacheMock) Has(key string, value interface{}) bool {
	if raw, found := c.values[key]; found {
		pattern, err := json.Marshal(value)
		if err != nil {
			return false
		}

		return string(pattern) == raw
	}

	return false
}
