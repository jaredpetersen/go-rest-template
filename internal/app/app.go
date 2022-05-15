package app

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jaredpetersen/go-health/health"
	"github.com/jaredpetersen/go-rest-template/api"
	"github.com/jaredpetersen/go-rest-template/internal/task"
	"github.com/rs/zerolog/log"
)

// Define interfaces where they are used

type TaskManager interface {
	Get(ctx context.Context, id string) (*task.Task, error)
	Save(ctx context.Context, t task.Task) error
}

type app struct {
	router        *chi.Mux
	HealthMonitor *health.Monitor
	TaskManager   TaskManager
}

type AppError struct {
	External error
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
	log.Error().AnErr("external", appErr.External).AnErr("internal", appErr.Internal).Send()

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if appErr.External != nil {
		json.NewEncoder(w).Encode(&api.Error{Message: appErr.External.Error()})
	}
}
