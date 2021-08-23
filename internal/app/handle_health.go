package app

import (
	"net/http"

	"github.com/jaredpetersen/go-rest-template/internal/health"
)

func (a *app) handleHealth() http.HandlerFunc {
	type response struct {
		Up bool `json:"up"`
	}

	return func(w http.ResponseWriter, req *http.Request) {
		status := health.Check()
		res := &response{
			Up: status.Up,
		}

		var statusCode int

		if status.Up {
			statusCode = http.StatusOK
		} else {
			statusCode = http.StatusServiceUnavailable
		}

		respond(w, res, statusCode)
	}
}
