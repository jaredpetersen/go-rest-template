package taskmgr

import (
	"context"

	"github.com/jaredpetersen/go-rest-template/internal/task"
	"github.com/rs/zerolog/log"
)

// DBRepo is a database repository for tasks.
type Manager struct {
	TaskCacheClient task.CacheClient
	TaskDBClient    task.DBClient
}

// Get retrieves a task by ID, first looking to the cache and then falling back on the database.
func (mgr Manager) Get(ctx context.Context, id string) (*task.Task, error) {
	task, err := mgr.TaskCacheClient.Get(ctx, id)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to retrieve task from cache")
	}

	if task != nil {
		return task, nil
	}

	return mgr.TaskDBClient.Get(ctx, id)
}

// Save stores a task to both cache and database.
//
// If the save to the cache fails, the error is logged and ignored so that we are resilient to fleeting cache
// dependency issues.
func (mgr Manager) Save(ctx context.Context, t task.Task) error {
	err := mgr.TaskCacheClient.Save(ctx, t)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to store task in cache")
	}

	return mgr.TaskDBClient.Save(ctx, t)
}
