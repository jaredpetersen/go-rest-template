package task

import (
	"context"
	"encoding/json"

	"github.com/jaredpetersen/go-rest-template/internal/redis"
)

type CacheClient interface {
	Get(ctx context.Context, id string) (*Task, error)
	Save(ctx context.Context, t Task) error
}

type CacheRepo struct {
	rdb redis.Client
}

func NewCacheRepo(rdb redis.Client) *CacheRepo {
	return &CacheRepo{rdb: rdb}
}

func (cr *CacheRepo) Get(ctx context.Context, id string) (*Task, error) {
	key := getRedisKey(id)
	val, err := cr.rdb.Get(ctx, key)
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

func (cr *CacheRepo) Save(ctx context.Context, t Task) error {
	key := getRedisKey(t.Id)
	value, err := json.Marshal(t)
	if err != nil {
		return err
	}

	return cr.rdb.Set(ctx, key, value, 0)
}

func getRedisKey(id string) string {
	return "task." + id
}