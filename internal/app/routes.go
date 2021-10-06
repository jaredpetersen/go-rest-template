package app

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
)

// routes sets up the server routing and configures the handlers
func (a *app) routes() {
	a.router = chi.NewRouter()

	// Set up logging middleware
	a.router.Use(hlog.NewHandler(log.Logger))
	a.router.Use(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		hlog.FromRequest(r).
			Info().
			Str("method", r.Method).
			Stringer("url", r.URL).
			Int("status", status).
			Int("size", size).
			Dur("duration", duration).
			Msg("Access")
	}))
	a.router.Use(hlog.RemoteAddrHandler("ip"))
	a.router.Use(hlog.UserAgentHandler("user_agent"))
	a.router.Use(hlog.RefererHandler("referer"))

	a.router.Get("/tasks/{id}", a.handleTaskGet())
	a.router.Post("/tasks", a.handleTaskSave())

	a.router.Get("/health", a.handleHealth())

	a.router.NotFound(a.handleNotFound())
	a.router.MethodNotAllowed(a.handleMethodNotAllowed())
}
