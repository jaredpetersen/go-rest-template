package server

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
)

// routes sets up the server routing and configures the handlers
func (srv *server) routes() {
	srv.router = mux.NewRouter()

	middlewareChain := newMiddlewareChain()

	srv.router.Handle("/health", middlewareChain.Then(srv.handleHealth())).Methods("GET")

	srv.router.NotFoundHandler = middlewareChain.Then(srv.handleNotFound())
	srv.router.MethodNotAllowedHandler = middlewareChain.Then(srv.handleMethodNotAllowed())
}

// newMiddlewareChain creates a general purpose middleware chain
func newMiddlewareChain() alice.Chain {
	middlewareChain := alice.New()
	middlewareChain = middlewareChain.Extend(newMiddlewareChainLog())

	return middlewareChain
}

// newMiddlewareChainLog creates a middleware chain for everything related to logging
func newMiddlewareChainLog() alice.Chain {
	middlewareChain := alice.New()
	middlewareChain = middlewareChain.Append(hlog.NewHandler(log.Logger))
	middlewareChain = middlewareChain.Append(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		hlog.FromRequest(r).
			Info().
			Str("method", r.Method).
			Stringer("url", r.URL).
			Int("status", status).
			Int("size", size).
			Dur("duration", duration).
			Msg("Access")
	}))
	middlewareChain = middlewareChain.Append(hlog.RemoteAddrHandler("ip"))
	middlewareChain = middlewareChain.Append(hlog.UserAgentHandler("user_agent"))
	middlewareChain = middlewareChain.Append(hlog.RefererHandler("referer"))

	return middlewareChain
}
