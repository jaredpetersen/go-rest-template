package app

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jaredpetersen/go-rest-template/internal/task"
)

func (a *app) handleTaskGet() http.HandlerFunc {
	type response struct {
		Id          string  `json:"id"`
		Description string  `json:"description"`
		DateDue     *string `json:"date_due"`
	}

	return func(w http.ResponseWriter, req *http.Request) {
		id := chi.URLParam(req, "id")

		val, err := task.Get(req.Context(), a.Redis, id)
		if err != nil {
			respondError(w, AppError{Internal: err}, http.StatusUnprocessableEntity)
			return
		}

		if val == nil {
			respond(w, nil, http.StatusNotFound)
			return
		}

		res := &response{
			Id:          val.Id,
			Description: val.Description,
			DateDue:     val.DateDue,
		}
		respond(w, res, http.StatusOK)
	}
}

func (a *app) handleTaskSave() http.HandlerFunc {
	type request struct {
		Description string  `json:"description"`
		DateDue     *string `json:"date_due"`
	}
	type response struct {
		Id string `json:"id"`
	}

	return func(w http.ResponseWriter, req *http.Request) {
		val := new(request)
		err := receive(req, val)
		if err != nil {
			respondError(w, AppError{Internal: err}, http.StatusBadRequest)
			return
		}

		t := task.New()
		t.Description = val.Description
		t.DateDue = val.DateDue

		err = task.Save(req.Context(), a.Redis, *t)
		if err != nil {
			respondError(w, AppError{Internal: err}, http.StatusUnprocessableEntity)
			return
		}

		respond(w, &response{t.Id}, http.StatusCreated)
	}
}
