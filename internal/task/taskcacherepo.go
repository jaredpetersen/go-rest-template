package task

import (
	"context"
	"encoding/json"

	"github.com/jaredpetersen/go-rest-template/internal/redis"
)

// CacheClient is a client for retrieving and manipulating tasks in the cache
type CacheClient interface {
	Get(ctx context.Context, id string) (*Task, error)
	Save(ctx context.Context, t Task) error
}

// CacheRepo is a cache repository for tasks.
type CacheRepo struct {
	Redis redis.Client
}

// Get retrieves a task from the cache using the task's ID. If a task cannot be found with that ID, nil will be
// returned for both the task and error.
func (cr CacheRepo) Get(ctx context.Context, id string) (*Task, error) {
	key := getRedisKey(id)
	val, err := cr.Redis.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, nil
	}

	var t Task
	err = json.Unmarshal([]byte(*val), &t)

	return &t, err
}

// Save stores a task in the cache.
func (cr CacheRepo) Save(ctx context.Context, t Task) error {
	key := getRedisKey(t.Id)
	value, err := json.Marshal(t)
	if err != nil {
		return err
	}

	return cr.Redis.Set(ctx, key, value, 0)
}

// getRedisKey builds a redis key for the task in the cache.
func getRedisKey(id string) string {
	return "task." + id
}
