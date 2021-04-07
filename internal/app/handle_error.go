package app

import "net/http"

func (a *app) handleNotFound() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(404)
	}
}

func (a *app) handleMethodNotAllowed() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(405)
	}
}
