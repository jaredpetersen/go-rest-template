package task

import (
	"context"
	"database/sql"
)

// Add more docs
// Integrate open API with codegen and only let it go as far as the handlers
// put open api in an api directory at the top level so that services can import your api
// via go get -- may require that we start checking in mocks

// DBClient is a client for retrieving and manipulating tasks in a SQL database
type DBClient interface {
	Get(ctx context.Context, id string) (*Task, error)
	Save(ctx context.Context, t Task) error
}

// DBRepo is a repository for tasks.
type DBRepo struct {
	DB sql.DB
}

// Get retrieves a task from the database using the task's ID. If a task cannot be found with that ID, nil will be
// returned for both the task and error.
func (dbr DBRepo) Get(ctx context.Context, id string) (*Task, error) {
	const query = `select description, date_due, date_created, date_updated
		from task
		where id = $1`
	row := dbr.DB.QueryRowContext(ctx, query, id)

	tsk := Task{Id: id}
	err := row.Scan(&tsk.Description, &tsk.DateDue, &tsk.DateCreated, &tsk.DateUpdated)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &tsk, nil
}

// Save stores a task in the database.
func (dbr DBRepo) Save(ctx context.Context, t Task) error {
	const query = `insert into "task" (id, description, date_due, date_created, date_updated)
		values ($1, $2, $3, $4, $5)`
	_, err := dbr.DB.ExecContext(ctx,
		query,
		t.Id,
		t.Description,
		t.DateDue,
		t.DateCreated,
		t.DateUpdated)

	return err
}
