package app

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jaredpetersen/go-rest-template/api"
	"github.com/jaredpetersen/go-rest-template/internal/task"
)

func (a *app) handleTaskGet() http.HandlerFunc {
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

		res := api.Task{
			Id:          val.ID,
			Description: val.Description,
			DateDue:     val.DateDue,
		}
		respond(w, res, http.StatusOK)
	}
}

func (a *app) handleTaskSave() http.HandlerFunc {
	// Set up dependencies specific to the handler here

	return func(w http.ResponseWriter, req *http.Request) {
		val := new(api.NewTask)
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

		respond(w, api.Identifier{Id: t.ID}, http.StatusCreated)
	}
}
