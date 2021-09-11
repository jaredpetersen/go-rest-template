package app

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/jaredpetersen/go-rest-template/internal/task"
	"github.com/jaredpetersen/go-rest-template/internal/tasksvc"
)

func (a *app) handleTaskGet() http.HandlerFunc {
	type response struct {
		Id          string     `json:"id"`
		Description string     `json:"description"`
		DateDue     *time.Time `json:"date_due,string"`
	}

	// Set up dependencies specific to the handler here
	tcr := task.CacheRepo{Redis: a.Redis}
	tdbr := task.DBRepo{DB: a.DB}

	return func(w http.ResponseWriter, req *http.Request) {
		id := chi.URLParam(req, "id")

		val, err := tasksvc.Get(req.Context(), tcr, tdbr, id)
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
		Description string     `json:"description"`
		DateDue     *time.Time `json:"date_due,string"`
	}
	type response struct {
		Id string `json:"id"`
	}

	// Set up dependencies specific to the handler here
	tcr := task.CacheRepo{Redis: a.Redis}
	tdbr := task.DBRepo{DB: a.DB}

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

		err = tasksvc.Save(req.Context(), tcr, tdbr, *t)
		if err != nil {
			respondError(w, AppError{Internal: err}, http.StatusUnprocessableEntity)
			return
		}

		respond(w, &response{t.Id}, http.StatusCreated)
	}
}
