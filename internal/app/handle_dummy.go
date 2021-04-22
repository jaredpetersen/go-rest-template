package app

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jaredpetersen/go-rest-example/internal/dummy"
)

func (a *app) handleDummyGet() http.HandlerFunc {
	type response struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}

	return func(w http.ResponseWriter, req *http.Request) {
		id := chi.URLParam(req, "id")

		val, err := dummy.Get(req.Context(), a.Redis, id)
		if err != nil {
			respondError(w, err.Error(), http.StatusUnprocessableEntity)
		}

		if val == nil {
			respondError(w, "", http.StatusNotFound)
		} else {
			res := &response{
				Id:   val.Id,
				Name: val.Name,
			}
			respond(w, res, http.StatusOK)
		}
	}
}

func (a *app) handleDummyNew() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		d := &dummy.Dummy{}
		err := receive(req, d)
		if err != nil {
			respondError(w, err.Error(), http.StatusUnprocessableEntity)
		}

		err = dummy.New(req.Context(), a.Redis, *d)
		if err != nil {
			respondError(w, err.Error(), http.StatusUnprocessableEntity)
		}
	}
}
