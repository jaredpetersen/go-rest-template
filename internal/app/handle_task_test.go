package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jaredpetersen/go-rest-template/internal/app/mocks"
	"github.com/jaredpetersen/go-rest-template/internal/task"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandleTaskGet(t *testing.T) {
	tsk := task.New()

	// Set up relevant server dependencies
	tskMgr := mocks.TaskManager{}
	tskMgr.On("Get", mock.Anything, tsk.Id).Return(tsk, nil)

	// Set up server
	a := New()
	a.TaskManager = &tskMgr

	// Make request
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/tasks/%s", tsk.Id), nil)
	require.NoError(t, err)
	res := httptest.NewRecorder()
	a.ServeHTTP(res, req)

	expectedJSON := fmt.Sprintf("{\"id\": \"%s\", \"description\": \"%s\", \"date_due\": null}", tsk.Id, tsk.Description)

	assert.Equal(t, http.StatusOK, res.Result().StatusCode)
	assert.JSONEq(t, expectedJSON, res.Body.String())
}

func TestHandleTaskGetError(t *testing.T) {
	// Set up relevant server dependencies
	tskMgr := mocks.TaskManager{}
	tskMgr.On("Get", mock.Anything, mock.AnythingOfType("string")).Return(nil, errors.New("failure to get task"))

	// Set up server
	a := New()
	a.TaskManager = &tskMgr

	// Make request
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/tasks/%s", uuid.New()), nil)
	require.NoError(t, err)
	res := httptest.NewRecorder()
	a.ServeHTTP(res, req)

	assert.Equal(t, http.StatusInternalServerError, res.Result().StatusCode)
	assert.Empty(t, res.Body)
}

func TestHandleTaskGetNotFound(t *testing.T) {
	// Set up relevant server dependencies
	tskMgr := mocks.TaskManager{}
	tskMgr.On("Get", mock.Anything, mock.AnythingOfType("string")).Return(nil, nil)

	// Set up server
	a := New()
	a.TaskManager = &tskMgr

	// Make request
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/tasks/%s", uuid.New().String()), nil)
	require.NoError(t, err)
	res := httptest.NewRecorder()
	a.ServeHTTP(res, req)

	assert.Equal(t, http.StatusNotFound, res.Result().StatusCode)
	assert.Empty(t, res.Body)
}

func TestHandleTaskSave(t *testing.T) {
	// Set up relevant server dependencies
	tskMgr := mocks.TaskManager{}
	tskMgr.On("Save", mock.Anything, mock.Anything).Return(nil)

	// Set up server
	a := New()
	a.TaskManager = &tskMgr

	// Set up request body
	tsk := struct {
		Description string `json:"description"`
	}{
		Description: "Buy milk",
	}
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(tsk)

	// Make request
	req, err := http.NewRequest(http.MethodPost, "/tasks", &buf)
	require.NoError(t, err)
	res := httptest.NewRecorder()
	a.ServeHTTP(res, req)

	assert.Equal(t, http.StatusCreated, res.Result().StatusCode)

	// Decode response body to struct so that we can pick out pieces
	resBody := struct {
		ID string `json:"id"`
	}{}
	err = json.NewDecoder(res.Body).Decode(&resBody)
	require.NoError(t, err, "Failed to convert response body")

	_, err = uuid.Parse(resBody.ID)
	require.NoError(t, err, "Returned an invalid UUID")
}

func TestHandleTaskSaveBadBody(t *testing.T) {
	// Set up relevant server dependencies
	tskMgr := mocks.TaskManager{}
	tskMgr.On("Save", mock.Anything, mock.Anything).Return(nil)

	// Set up server
	a := New()
	a.TaskManager = &tskMgr

	// Set up request body without valid JSON
	reqBody := strings.NewReader("<task />")

	// Make request
	req, err := http.NewRequest(http.MethodPost, "/tasks", reqBody)
	require.NoError(t, err)
	res := httptest.NewRecorder()
	a.ServeHTTP(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode)
	assert.Empty(t, res.Body)
}

func TestHandleTaskSaveMissingBodyFields(t *testing.T) {
	// Set up relevant server dependencies
	tskMgr := mocks.TaskManager{}
	tskMgr.On("Save", mock.Anything, mock.Anything).Return(nil)

	// Set up server
	a := New()
	a.TaskManager = &tskMgr

	// Set up invalid request body
	tsk := struct {
		FavoriteColor string `json:"favorite_color"`
	}{
		FavoriteColor: "Seafoam Green",
	}
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(tsk)

	// Make request
	req, err := http.NewRequest(http.MethodPost, "/tasks", &buf)
	require.NoError(t, err)
	res := httptest.NewRecorder()
	a.ServeHTTP(res, req)

	assert.Equal(t, http.StatusUnprocessableEntity, res.Result().StatusCode)
	assert.JSONEq(t, "{\"message\": \"Field 'description' is required\"}", res.Body.String())
}

func TestHandleTaskSaveError(t *testing.T) {
	// Set up relevant server dependencies
	tskMgr := mocks.TaskManager{}
	tskMgr.On("Save", mock.Anything, mock.Anything).Return(errors.New("failure to save task"))

	// Set up server
	a := New()
	a.TaskManager = &tskMgr

	// Set up request body
	tsk := struct {
		Description string `json:"description"`
	}{
		Description: "Buy milk",
	}
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(tsk)

	// Make request
	req, err := http.NewRequest(http.MethodPost, "/tasks", &buf)
	require.NoError(t, err)
	res := httptest.NewRecorder()
	a.ServeHTTP(res, req)

	assert.Equal(t, http.StatusInternalServerError, res.Result().StatusCode)
	assert.Empty(t, res.Body)
}
