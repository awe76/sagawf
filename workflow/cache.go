package workflow

import (
	"context"
	"encoding/json"
	"fmt"

	mc "go-micro.dev/v4/cache"
)

type cache struct {
	cache mc.Cache
}

func NewCache(opts ...mc.Option) Cache {
	c := mc.NewCache(opts...)
	return &cache{c}
}

type Cache interface {
	Set(ctx context.Context, key string, value interface{}) error
	Get(ctx context.Context, key string) (string, error)
	Remove(ctx context.Context, key string) error
}

func (c *cache) Set(ctx context.Context, key string, value interface{}) error {
	rawValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.cache.Context(ctx).Put(key, string(rawValue), 0)
}

func (c *cache) Get(ctx context.Context, key string) (string, error) {
	value, _, err := c.cache.Context(ctx).Get(key)
	if err != nil {
		return "", err
	}

	result := fmt.Sprintf("%v", value)
	return result, nil
}

func (c *cache) Remove(ctx context.Context, key string) error {

	return c.cache.Context(ctx).Delete(key)
}
