package task

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
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
			date_due timestamp,
			date_created timestamp not null,
			date_updated timestamp not null);`
	_, err := db.ExecContext(ctx, query)

	return err
}

func TestDBRepoSaveGet(t *testing.T) {
	ctx := context.Background()

	cdbContainer, err := setupCockroachDB(ctx)
	require.NoError(t, err, "Failed to start up CockroachDB container")
	defer cdbContainer.Terminate(ctx)

	db, err := sql.Open("pgx", cdbContainer.URI+"/projectmanagement")
	require.NoError(t, err, "Failed to open connection to CockroachDB")
	defer db.Close()

	err = initCockroachDB(ctx, *db)

	tsk := New()

	tdbr := NewDBRepo(*db)

	// TODO table-driven tests

	err = tdbr.Save(ctx, *tsk)
	require.NoError(t, err, "Save returned error")

	savedTsk, err := tdbr.Get(ctx, tsk.Id)
	require.NoError(t, err, "Get returned error")
	require.NotNil(t, savedTsk)
	assert.True(t, cmp.Equal(*tsk, *savedTsk), "Saved task is not the same:\n"+cmp.Diff(*tsk, *savedTsk))
}
