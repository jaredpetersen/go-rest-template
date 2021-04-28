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
			respondError(w, AppError{Internal: err}, http.StatusUnprocessableEntity)
			return
		}

		if val == nil {
			respond(w, nil, http.StatusNotFound)
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
	type request struct {
		Name string `json:"name"`
	}

	return func(w http.ResponseWriter, req *http.Request) {
		val := request{}
		err := receive(req, val)
		if err != nil {
			respondError(w, AppError{Internal: err}, http.StatusBadRequest)
			return
		}

		err = dummy.New(req.Context(), a.Redis, (*dummy.Dummy)(&val))

		if err != nil {
			respondError(w, AppError{Internal: err}, http.StatusUnprocessableEntity)
			return
		}
	}
}
