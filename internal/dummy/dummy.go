package dummy

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type Dummy struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	DateCreated string `json:"date_created"`
	LastUpdated string `json:"last_updated"`
}

func Get(ctx context.Context, rdb *redis.Client, id string) (*Dummy, error) {
	dummyBytes, err := getCached(ctx, rdb, id)
	if err != nil {
		log.Error().Err(err).Msg("get")
		return nil, err
	}
	if dummyBytes == nil {
		return nil, nil
	}

	var dummy Dummy
	json.Unmarshal(dummyBytes, &dummy)

	return &dummy, nil
}

func New(ctx context.Context, rdb *redis.Client, dummy Dummy) error {
	if dummy.Id == "" {
		dummy.Id = uuid.New().String()
	}
	if dummy.DateCreated == "" {
		dummy.DateCreated = time.Now().Format(time.RFC3339)
	}
	if dummy.LastUpdated == "" {
		dummy.LastUpdated = time.Now().Format(time.RFC3339)
	}

	key := "DUMMY." + dummy.Id
	value, err := json.Marshal(dummy)
	if err != nil {
		return err
	}

	return rdb.Set(ctx, key, value, 0).Err()
}

func getCached(ctx context.Context, rdb *redis.Client, id string) ([]byte, error) {
	key := "DUMMY." + id
	val, err := rdb.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return val, nil
}
