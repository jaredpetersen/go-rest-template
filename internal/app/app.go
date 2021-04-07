package app

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
)

type app struct {
	router *chi.Mux
	// TODO DB
	// TODO Redis
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

func respondError(w http.ResponseWriter, error string, statusCode int) {
	type response struct {
		Message string `json:"message"`
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(&response{Message: error})
}
