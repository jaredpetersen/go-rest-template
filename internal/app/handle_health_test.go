package app

import (
	"context"
	"encoding/json"
	"github.com/jaredpetersen/go-health/health"
	"github.com/jaredpetersen/go-rest-template/api"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func buildHealthCheckFunc(status health.Status) health.CheckFunc {
	return func(ctx context.Context) health.Status {
		return status
	}
}

func TestHandleLiveness(t *testing.T) {
	// Set up server
	a := New()

	// Make request
	req, err := http.NewRequest(http.MethodGet, "/liveness", nil)
	require.NoError(t, err)
	res := httptest.NewRecorder()
	a.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Result().StatusCode)
	assert.Empty(t, res.Body)
}

func TestHandleReadinessStateUp(t *testing.T) {
	ctx := context.Background()

	// Set up health monitor and wait for it to kick off monitoring goroutines
	dbHealthCheck := health.NewCheck("database", buildHealthCheckFunc(health.Status{State: health.StateUp}))
	redisHealthCheck := health.NewCheck("redis", buildHealthCheckFunc(health.Status{State: health.StateUp}))

	healthMonitor := health.New()
	healthMonitor.Monitor(ctx, redisHealthCheck, dbHealthCheck)
	time.Sleep(time.Millisecond * 200)

	// Set up server
	a := New()
	a.HealthMonitor = healthMonitor

	// Make request
	req, err := http.NewRequest(http.MethodGet, "/readiness", nil)
	require.NoError(t, err)
	res := httptest.NewRecorder()
	a.ServeHTTP(res, req)

	// Decode response body to struct so that we can pick out pieces
	resBody := api.Health{}
	err = json.NewDecoder(res.Body).Decode(&resBody)
	require.NoError(t, err, "Failed to convert response body")

	assert.Equal(t, http.StatusOK, res.Result().StatusCode)
	assert.Equal(t, api.HealthStateUP, resBody.State)
	assert.Equal(t, api.HealthStateUP, resBody.Components.CockroachDb.State)
	assert.NotNil(t, resBody.Components.CockroachDb.Timestamp)
	assert.Equal(t, api.HealthStateUP, resBody.Components.Redis.State)
	assert.NotNil(t, resBody.Components.Redis.Timestamp)
}

func TestHandleReadinessStateWarn(t *testing.T) {
	ctx := context.Background()

	// Set up health monitor and wait for it to kick off monitoring goroutines
	dbHealthCheck := health.NewCheck("database", buildHealthCheckFunc(health.Status{State: health.StateUp}))
	redisHealthCheck := health.NewCheck("redis", buildHealthCheckFunc(health.Status{State: health.StateWarn}))

	healthMonitor := health.New()
	healthMonitor.Monitor(ctx, redisHealthCheck, dbHealthCheck)
	time.Sleep(time.Millisecond * 200)

	// Set up server
	a := New()
	a.HealthMonitor = healthMonitor

	// Make request
	req, err := http.NewRequest(http.MethodGet, "/readiness", nil)
	require.NoError(t, err)
	res := httptest.NewRecorder()
	a.ServeHTTP(res, req)

	// Decode response body to struct so that we can pick out pieces
	resBody := api.Health{}
	err = json.NewDecoder(res.Body).Decode(&resBody)
	require.NoError(t, err, "Failed to convert response body")

	assert.Equal(t, http.StatusOK, res.Result().StatusCode)
	assert.Equal(t, api.HealthStateWARN, resBody.State)
	assert.Equal(t, api.HealthStateUP, resBody.Components.CockroachDb.State)
	assert.NotNil(t, resBody.Components.CockroachDb.Timestamp)
	assert.Equal(t, api.HealthStateWARN, resBody.Components.Redis.State)
	assert.NotNil(t, resBody.Components.Redis.Timestamp)
}

func TestHandleReadinessStateDown(t *testing.T) {
	ctx := context.Background()

	// Set up health monitor and wait for it to kick off monitoring goroutines
	dbHealthCheck := health.NewCheck("database", buildHealthCheckFunc(health.Status{State: health.StateDown}))
	redisHealthCheck := health.NewCheck("redis", buildHealthCheckFunc(health.Status{State: health.StateUp}))

	healthMonitor := health.New()
	healthMonitor.Monitor(ctx, redisHealthCheck, dbHealthCheck)
	time.Sleep(time.Millisecond * 200)

	// Set up server
	a := New()
	a.HealthMonitor = healthMonitor

	// Make request
	req, err := http.NewRequest(http.MethodGet, "/readiness", nil)
	require.NoError(t, err)
	res := httptest.NewRecorder()
	a.ServeHTTP(res, req)

	// Decode response body to struct so that we can pick out pieces
	resBody := api.Health{}
	err = json.NewDecoder(res.Body).Decode(&resBody)
	require.NoError(t, err, "Failed to convert response body")

	assert.Equal(t, http.StatusServiceUnavailable, res.Result().StatusCode)
	assert.Equal(t, api.HealthStateDOWN, resBody.State)
	assert.Equal(t, api.HealthStateDOWN, resBody.Components.CockroachDb.State)
	assert.NotNil(t, resBody.Components.CockroachDb.Timestamp)
	assert.Equal(t, api.HealthStateUP, resBody.Components.Redis.State)
	assert.NotNil(t, resBody.Components.Redis.Timestamp)
}
