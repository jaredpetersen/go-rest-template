package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jaredpetersen/go-rest-template/internal/app"
	"github.com/jaredpetersen/go-rest-template/internal/redis"
	"github.com/jaredpetersen/go-rest-template/internal/task"
	"github.com/jaredpetersen/go-rest-template/internal/taskmgr"
	"github.com/rs/zerolog/log"
)

func main() {
	// TODO app config (maybe via kelseyhightower/envconfig)

	a := app.New()

	// Set up SQL database
	db, err := sql.Open("pgx", "postgres://go-api:password@localhost:26257/project-management")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer db.Close()

	// Set up Redis
	rCfg := redis.Config{URI: "redis://localhost:6379"}
	rdb, err := redis.New(rCfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to Redis")
	}
	defer rdb.Close()

	// Set up task manager
	taskCacheClient := task.CacheRepo{Redis: rdb}
	taskDBClient := task.DBRepo{DB: *db}
	a.TaskManager = taskmgr.Manager{TaskDBClient: taskDBClient, TaskCacheClient: taskCacheClient}

	addr := 8080
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", addr),
		Handler:      a,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Info().Int("port", addr).Msg("Started")

	err = srv.ListenAndServe()
	if err != nil {
		log.Error().Err(err).Msg("Server encountered an error")
	}
}
