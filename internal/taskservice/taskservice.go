package taskservice

import (
	"context"

	"github.com/jaredpetersen/go-rest-template/internal/task"
	"github.com/rs/zerolog/log"
)

// Get retrieves a task by ID, first looking to the cache and then falling back on the database.
func Get(ctx context.Context, tcr task.CacheClient, tdbr task.DBClient, id string) (*task.Task, error) {
	task, err := tcr.Get(ctx, id)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to retrieve task from cache")
	}

	if task != nil {
		return task, nil
	}

	return tdbr.Get(ctx, id)
}

// Save stores a task to both cache and database.
//
// If the save to the cache fails, the error is logged and ignored so that we are resilient to fleeting cache
// dependency issues.
func Save(ctx context.Context, tcr task.CacheClient, tdbr task.DBClient, t task.Task) error {
	err := tcr.Save(ctx, t)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to store task in cache")
	}

	return tdbr.Save(ctx, t)
}
