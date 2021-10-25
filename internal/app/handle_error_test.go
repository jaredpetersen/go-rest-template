package app_test

import (
	"github.com/jaredpetersen/go-rest-template/internal/app"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleNotFound(t *testing.T) {
	// Set up server
	a := app.New()

	// Make request
	req, err := http.NewRequest(http.MethodGet, "/bleepbloop", nil)
	require.NoError(t, err)
	res := httptest.NewRecorder()
	a.ServeHTTP(res, req)

	assert.Equal(t, http.StatusNotFound, res.Result().StatusCode)
	assert.Empty(t, res.Body)
}

func TestHandleMethodNotAllowed(t *testing.T) {
	// Set up server
	a := app.New()

	// Make request
	req, err := http.NewRequest(http.MethodPost, "/readiness", nil)
	require.NoError(t, err)
	res := httptest.NewRecorder()
	a.ServeHTTP(res, req)

	assert.Equal(t, http.StatusMethodNotAllowed, res.Result().StatusCode)
	assert.Empty(t, res.Body)
}
