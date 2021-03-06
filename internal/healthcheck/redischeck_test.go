package healthcheck_test

import (
	"context"
	"errors"
	"github.com/jaredpetersen/go-health/health"
	"github.com/jaredpetersen/go-rest-template/internal/healthcheck"
	redismock "github.com/jaredpetersen/go-rest-template/internal/redis/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBuildRedisHealthCheckFuncStateWarn(t *testing.T) {
	ctx := context.Background()

	rdb := redismock.Client{}
	rdb.On("Ping", ctx).Return(errors.New("bad ping"))

	redisHealthCheckFunc := healthcheck.BuildRedisHealthCheckFunc(&rdb)
	require.NotNil(t, redisHealthCheckFunc)

	healthStatus := redisHealthCheckFunc(ctx)
	assert.Equal(t, health.StateWarn, healthStatus.State)
	assert.Nil(t, healthStatus.Details)
}

func TestBuildRedisHealthCheckFuncStateUp(t *testing.T) {
	ctx := context.Background()

	rdb := redismock.Client{}
	rdb.On("Ping", ctx).Return(nil)

	redisHealthCheckFunc := healthcheck.BuildRedisHealthCheckFunc(&rdb)
	require.NotNil(t, redisHealthCheckFunc)

	healthStatus := redisHealthCheckFunc(ctx)
	assert.Equal(t, health.StateUp, healthStatus.State)
	assert.Nil(t, healthStatus.Details)
}
