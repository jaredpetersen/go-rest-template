package task

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
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

func initCockroachDB(ctx context.Context, db sql.DB) error {
	const query = `CREATE DATABASE projectmanagement;
		CREATE TABLE projectmanagement.task(
			id uuid primary key not null,
			description varchar(255) not null,
			date_due timestamp with time zone,
			date_created timestamp with time zone not null,
			date_updated timestamp with time zone not null);`
	_, err := db.ExecContext(ctx, query)

	return err
}

func truncateCockroachDB(ctx context.Context, db sql.DB) error {
	const query = `truncate projectmanagement.task`
	_, err := db.ExecContext(ctx, query)
	return err
}

func TestIntegrationDBRepoSaveGet(t *testing.T) {
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

	err = initCockroachDB(ctx, *db)
	require.NoError(t, err, "Failed to initialize CockroachDB")
	defer truncateCockroachDB(ctx, *db)

	tdbr := NewDBRepo(*db)

	now := time.Now()

	var tests = []*Task{
		New(),
		func() *Task {
			tsk := New()
			tsk.Description = "Update resumé"
			return tsk
		}(),
		func() *Task {
			tsk := New()
			tsk.Description = "Call veterinarian"
			tsk.DateDue = &now
			return tsk
		}(),
	}

	for _, tt := range tests {
		defer truncateCockroachDB(ctx, *db)

		err = tdbr.Save(ctx, *tt)
		require.NoError(t, err, "Save returned error")

		savedTsk, err := tdbr.Get(ctx, tt.Id)
		require.NoError(t, err, "Get returned error")
		require.NotNil(t, savedTsk, "Get did not return a task")
		assert.True(t, cmp.Equal(*tt, *savedTsk), "Saved task is not the same:\n"+cmp.Diff(*tt, *savedTsk))
	}
}

func TestIntegrationDBRepoSaveDBError(t *testing.T) {
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

	// Do not initialize database tables

	tdbr := NewDBRepo(*db)

	err = tdbr.Save(ctx, *New())
	require.Error(t, err, "Save did not return error")
}

func TestIntegrationDBRepoGetNonexistent(t *testing.T) {
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

	err = initCockroachDB(ctx, *db)
	require.NoError(t, err, "Failed to initialize CockroachDB")
	defer truncateCockroachDB(ctx, *db)

	tdbr := NewDBRepo(*db)

	tsk, err := tdbr.Get(ctx, uuid.NewString())
	require.NoError(t, err, "Get returned error")
	assert.Nil(t, tsk, "Get returned a task")
}

func TestIntegrationDBRepoGetDBError(t *testing.T) {
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

	// Do not initialize database tables

	tdbr := NewDBRepo(*db)

	tsk, err := tdbr.Get(ctx, uuid.NewString())
	require.Error(t, err, "Get did not return error")
	assert.Nil(t, tsk, "Get returned a task")
}
