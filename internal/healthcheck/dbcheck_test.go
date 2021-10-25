package healthcheck

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jaredpetersen/go-health/health"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
)

type cockroachDBContainer struct {
	testcontainers.Container
	URI string
}

func setupCockroachDB(ctx context.Context) (*cockroachDBContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        "cockroachdb/cockroach:latest-v21.1",
		ExposedPorts: []string{"26257/tcp", "8080/tcp"},
		WaitingFor:   wait.ForHTTP("/health").WithPort("8080"),
		Cmd:          []string{"start-single-node", "--insecure"},
		SkipReaper:   true,
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, "26257")
	if err != nil {
		return nil, err
	}

	hostIP, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("postgres://root@%s:%s", hostIP, mappedPort.Port())

	return &cockroachDBContainer{Container: container, URI: uri}, nil
}

func TestBuildDBHealthCheckFuncStateDown(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	cdbContainer, err := setupCockroachDB(ctx)
	require.NoError(t, err, "Failed to start up CockroachDB container")
	defer cdbContainer.Terminate(ctx)

	db, err := sql.Open("pgx", cdbContainer.URI+"/projectmanagement")
	require.NoError(t, err, "Failed to open connection to CockroachDB")
	defer db.Close()

	dbHealthCheckFunc := BuildDBHealthCheckFunc(db)
	require.NotNil(t, dbHealthCheckFunc)

	// Simulate a database outage
	cdbContainer.Terminate(ctx)

	healthStatus := dbHealthCheckFunc(ctx)
	assert.Equal(t, health.StateDown, healthStatus.State)
	assert.Nil(t, healthStatus.Details)
}

func TestBuildDBHealthCheckFuncStateUp(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	cdbContainer, err := setupCockroachDB(ctx)
	require.NoError(t, err, "Failed to start up CockroachDB container")
	defer cdbContainer.Terminate(ctx)

	db, err := sql.Open("pgx", cdbContainer.URI+"/projectmanagement")
	require.NoError(t, err, "Failed to open connection to CockroachDB")
	defer db.Close()

	dbHealthCheckFunc := BuildDBHealthCheckFunc(db)
	require.NotNil(t, dbHealthCheckFunc)

	healthStatus := dbHealthCheckFunc(ctx)
	assert.Equal(t, health.StateUp, healthStatus.State)
	assert.Equal(t, DBDetails{ConnectionsInUse: 0, ConnectionsIdle: 1}, healthStatus.Details)
}