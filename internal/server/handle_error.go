package server

import "net/http"

func (srv *server) handleNotFound() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(404)
	}
}

func (srv *server) handleMethodNotAllowed() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(405)
	}
}
