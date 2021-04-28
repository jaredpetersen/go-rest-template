package app

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

type app struct {
	router *chi.Mux
	// TODO DB
	Redis *redis.Client
}

type AppError struct {
	Public   error
	Internal error
}

func New() *app {
	a := &app{}
	a.routes()
	return a
}

// ServeHTTP turns the app into an HTTP handler
func (a *app) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	a.router.ServeHTTP(w, req)
}

func receive(req *http.Request, data interface{}) error {
	// raw, _ := io.ReadAll(req.Body)
	// log.Debug().Str("raw", string(raw)).Send()
	return json.NewDecoder(req.Body).Decode(data)
}

func respond(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func respondError(w http.ResponseWriter, appErr AppError, statusCode int) {
	type response struct {
		Message string `json:"message"`
	}

	log.Error().AnErr("public", appErr.Public).AnErr("internal", appErr.Internal).Send()

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if appErr.Public != nil {
		json.NewEncoder(w).Encode(&response{Message: appErr.Public.Error()})
	}
}
