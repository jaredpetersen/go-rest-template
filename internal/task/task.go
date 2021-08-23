package task

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jaredpetersen/go-rest-example/internal/redis"
	"github.com/rs/zerolog/log"
)

// Task represents something that must be done
type Task struct {
	Id          string  `json:"id"`
	Description string  `json:"description"`
	DateDue     *string `json:"date_due"`
	DateCreated string  `json:"date_created"`
	DateUpdated string  `json:"date_updated"`
}

// New creates a new task with default values
func New() *Task {
	now := time.Now().Format(time.RFC3339)
	return &Task{Id: uuid.New().String(), DateCreated: now, DateUpdated: now}
}

// Get retrieves a task by ID
func Get(ctx context.Context, rdb redis.Client, id string) (*Task, error) {
	task, err := getCache(ctx, rdb, id)
	if err != nil {
		log.Error().Err(err).Msg("get")
		return nil, err
	}

	return task, nil
}

// Save stores a task
func Save(ctx context.Context, rdb redis.Client, t Task) error {
	return setCache(ctx, rdb, t)
}

func getRedisKey(id string) string {
	return "task." + id
}

func getCache(ctx context.Context, rdb redis.Client, id string) (*Task, error) {
	key := getRedisKey(id)
	val, err := rdb.Get(ctx, key)
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

func setCache(ctx context.Context, rdb redis.Client, t Task) error {
	key := getRedisKey(t.Id)
	value, err := json.Marshal(t)
	if err != nil {
		return err
	}

	return rdb.Set(ctx, key, value, 0)
}
