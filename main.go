package main

import (
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jaredpetersen/go-rest-example/internal/app"
	"github.com/rs/zerolog/log"
)

func main() {
	a := app.New()

	// Setup Redis
	redisOptions, err := redis.ParseURL("redis://localhost:6379")
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to Redis")
		os.Exit(1)
	}

	a.Redis = redis.NewClient(redisOptions)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      a,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Error().Err(err).Msg("Server encountered an error")
	}
}
