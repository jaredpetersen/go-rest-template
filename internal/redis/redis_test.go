package redis_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jaredpetersen/go-rest-template/internal/redis"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type redisContainer struct {
	testcontainers.Container
	URI string
}

// dummyDocument is a generic, bland data object that is used to verify that structs get turned into json strings
type dummyDocument struct {
	Name      string `json:"name"`
	LastCount int    `json:"last_count"`
	Valid     bool   `json:"valid"`
}

func (d dummyDocument) MarshalBinary() (data []byte, err error) {
	return json.Marshal(d)
}

func (d dummyDocument) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &d)
}

// setupRedis starts up a Redis container
//
// Returned Redis container must be explicitly terminated
func setupRedis(ctx context.Context) (*redisContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        "redis:6",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("* Ready to accept connections"),
		SkipReaper:   true,
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, "6379")
	if err != nil {
		return nil, err
	}

	hostIP, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("redis://%s:%s", hostIP, mappedPort.Port())

	return &redisContainer{Container: container, URI: uri}, nil
}

func TestNew(t *testing.T) {
	config := redis.Config{URI: "redis://localhost:6379"}

	rdb, err := redis.New(config)
	require.NoError(t, err, "Returned error")
	if assert.NotNil(t, rdb, "Client is nil") {
		defer rdb.Close()
	}
}

func TestNewReturnsConfigError(t *testing.T) {
	config := redis.Config{URI: "invaliduri"}

	rdb, err := redis.New(config)
	assert.Error(t, err, "Returned error")
	if !assert.Nil(t, rdb, "Client is nil") {
		defer rdb.Close()
	}
}

func TestIntegrationPing(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	redisContainer, err := setupRedis(ctx)
	require.NoError(t, err, "Failed to start up Redis container")
	defer redisContainer.Terminate(ctx)

	config := redis.Config{URI: redisContainer.URI}
	rdb, err := redis.New(config)
	require.NoError(t, err, "Client instantiation error")
	defer rdb.Close()

	err = rdb.Ping(ctx)
	assert.NoError(t, err, "Ping error")
}

func TestIntegrationGetNoKey(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	redisContainer, err := setupRedis(ctx)
	require.NoError(t, err, "Failed to start up Redis container")
	defer redisContainer.Terminate(ctx)

	config := redis.Config{URI: redisContainer.URI}
	rdb, err := redis.New(config)
	require.NoError(t, err, "Client instantiation error")
	defer rdb.Close()

	val, err := rdb.Get(ctx, "doesnotexist")
	assert.NoError(t, err, "Get error")

	assert.Nil(t, val, "Retrieved value pointer is not nil")
}

func TestIntegrationSetGet(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	redisContainer, err := setupRedis(ctx)
	require.NoError(t, err, "Failed to start up Redis container")
	defer redisContainer.Terminate(ctx)

	config := redis.Config{URI: redisContainer.URI}
	rdb, err := redis.New(config)
	require.NoError(t, err, "Client instantiation error")
	defer rdb.Close()

	var tests = []struct {
		input    interface{}
		expected string
	}{
		{
			input:    "turkey",
			expected: "turkey",
		},
		{
			input:    94,
			expected: "94",
		},
		{
			input:    false,
			expected: "0",
		},
		{
			input:    dummyDocument{Name: "john", LastCount: 5, Valid: true},
			expected: "{\"name\":\"john\",\"last_count\":5,\"valid\":true}",
		},
	}

	for _, tt := range tests {
		key := "dummy." + uuid.NewString()

		err = rdb.Set(ctx, key, tt.input, 0)
		assert.NoError(t, err, "Set error")

		val, err := rdb.Get(ctx, key)
		assert.NoError(t, err, "Get error")

		assert.Equal(t, tt.expected, *val, "Retrieved value is not equal")
	}
}

func TestIntegrationSetTTL(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	redisContainer, err := setupRedis(ctx)
	require.NoError(t, err, "Failed to start up Redis container")
	defer redisContainer.Terminate(ctx)

	config := redis.Config{URI: redisContainer.URI}
	rdb, err := redis.New(config)
	require.NoError(t, err, "Client instantiation error")
	defer rdb.Close()

	key := "dummy." + uuid.NewString()
	exp, _ := time.ParseDuration("30m")
	err = rdb.Set(ctx, key, "expiring", exp)
	assert.NoError(t, err, "Get error")

	ttl, err := rdb.TTL(ctx, key)
	assert.NoError(t, err, "TTL error")
	t.Logf("%s", ttl.String())
	assert.LessOrEqual(t, ttl.Seconds(), exp.Seconds())
	assert.GreaterOrEqual(t, ttl.Seconds(), 20.0)
}

func TestIntegrationSetTTLNoExpiration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	redisContainer, err := setupRedis(ctx)
	require.NoError(t, err, "Failed to start up Redis container")
	defer redisContainer.Terminate(ctx)

	config := redis.Config{URI: redisContainer.URI}
	rdb, err := redis.New(config)
	require.NoError(t, err, "Client instantiation error")
	defer rdb.Close()

	key := "dummy." + uuid.NewString()
	err = rdb.Set(ctx, key, "expiring", 0)
	assert.NoError(t, err, "Get error")

	ttl, err := rdb.TTL(ctx, key)
	assert.NoError(t, err, "TTL error")
	assert.Equal(t, time.Duration(-1), ttl)
}

func TestIntegrationTTLNoKey(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	redisContainer, err := setupRedis(ctx)
	require.NoError(t, err, "Failed to start up Redis container")
	defer redisContainer.Terminate(ctx)

	config := redis.Config{URI: redisContainer.URI}
	rdb, err := redis.New(config)
	require.NoError(t, err, "Client instantiation error")
	defer rdb.Close()

	ttl, err := rdb.TTL(ctx, "dummy")
	assert.NoError(t, err, "TTL error")
	assert.Equal(t, time.Duration(-2), ttl)
}

func TestIntegrationCloset(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	redisContainer, err := setupRedis(ctx)
	require.NoError(t, err, "Failed to start up Redis container")
	defer redisContainer.Terminate(ctx)

	config := redis.Config{URI: redisContainer.URI}
	rdb, err := redis.New(config)
	require.NoError(t, err, "Client instantiation error")
	defer rdb.Close()
}
