package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jaredpetersen/go-rest-example/internal/app"
	"github.com/jaredpetersen/go-rest-example/internal/redis"
	"github.com/rs/zerolog/log"
)

func main() {
	a := app.New()

	// Setup Redis
	rCfg := redis.Config{URI: "redis://localhost:6379"}
	rdb, err := redis.New(rCfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to Redis")
	}
	a.Redis = rdb

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
