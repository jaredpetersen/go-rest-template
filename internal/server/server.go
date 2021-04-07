package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type server struct {
	router *mux.Router
	// TODO DB
	// TODO Redis
}

func New() *server {
	srv := &server{}
	srv.routes()
	return srv
}

// ServeHTTP turns the server into an HTTP handler
func (srv *server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	srv.router.ServeHTTP(w, req)
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
