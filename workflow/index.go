package workflow

import (
	"context"
	"encoding/json"

	mc "go-micro.dev/v4/cache"
)

type Index struct {
	ID int
}

func ReserveID(key string, cache Cache) (int, error) {
	var index Index

	ctx := context.Background()
	rawIndex, err := cache.Get(ctx, key)

	if err == mc.ErrKeyNotFound {
		index = Index{
			ID: 0,
		}
	} else if err != nil {
		return index.ID, err
	} else {
		json.Unmarshal([]byte(rawIndex), &index)
	}

	index.ID++

	cache.Set(ctx, key, index)
	return index.ID, nil
}
