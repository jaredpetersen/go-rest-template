package app

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/jaredpetersen/go-rest-template/internal/task"
)

func (a *app) handleTaskGet() http.HandlerFunc {
	type response struct {
		Id          string     `json:"id"`
		Description string     `json:"description"`
		DateDue     *time.Time `json:"date_due,string"`
	}

	// Set up any dependencies specific to the handler here

	return func(w http.ResponseWriter, req *http.Request) {
		id := chi.URLParam(req, "id")

		val, err := a.TaskManager.Get(req.Context(), id)
		if err != nil {
			respondError(w, AppError{Internal: err}, http.StatusInternalServerError)
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
		Description string     `json:"description"`
		DateDue     *time.Time `json:"date_due,string"`
	}
	type response struct {
		Id string `json:"id"`
	}

	// Set up dependencies specific to the handler here

	return func(w http.ResponseWriter, req *http.Request) {
		val := new(request)
		err := receive(req, val)
		if err != nil {
			respondError(w, AppError{Internal: err}, http.StatusBadRequest)
			return
		}

		// Validate request body manually
		if val.Description == "" {
			respondError(w, AppError{External: errors.New("Field 'description' is required")}, http.StatusUnprocessableEntity)
			return
		}

		t := task.New()
		t.Description = val.Description
		t.DateDue = val.DateDue

		err = a.TaskManager.Save(req.Context(), *t)
		if err != nil {
			respondError(w, AppError{Internal: err}, http.StatusInternalServerError)
			return
		}

		respond(w, &response{t.Id}, http.StatusCreated)
	}
}
