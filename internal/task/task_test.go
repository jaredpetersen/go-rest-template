package task

import (
	"context"
	"testing"
	"time"

	"github.com/jaredpetersen/go-rest-example/internal/redis"
)

type redisClientStub struct {
	redis.Client
	getFunc func(ctx context.Context, key string) (*string, error)
	setFunc func(ctx context.Context, key string, value interface{}, expiration time.Duration) error
}

func (rs *redisClientStub) Get(ctx context.Context, key string) (*string, error) {
	return rs.getFunc(ctx, key)
}

func (rs *redisClientStub) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return rs.setFunc(ctx, key, value, expiration)
}

func TestSave(t *testing.T) {
	ctx := context.Background()

	rdb := redisClientStub{}
	rdb.setFunc = func(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
		return nil
	}

	task := Task{}

	err := Save(ctx, &rdb, task)
	if err != nil {
		t.Error("Encountered error", err)
	}
}

func TestGet(t *testing.T) {
	ctx := context.Background()

	rdb := redisClientStub{}
	rdb.getFunc = func(ctx context.Context, key string) (s *string, e error) {
		val := "{\"description\":\"buy socks\"}"
		return &val, nil
	}

	id := "2b7e1292-a831-4df5-b00e-3105a51111bb"

	_, err := Get(ctx, &rdb, id)
	if err != nil {
		t.Error("Encountered error", err)
	}
}
