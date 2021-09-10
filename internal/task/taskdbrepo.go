package task

import (
	"context"
	"database/sql"
)

type DBClient interface {
	Get(ctx context.Context, id string) (*Task, error)
	Save(ctx context.Context, t Task) error
}

type DBRepo struct {
	sdb sql.DB
}

func NewDBRepo(sdb sql.DB) *DBRepo {
	return &DBRepo{sdb: sdb}
}

func (dbr *DBRepo) Get(ctx context.Context, id string) (*Task, error) {
	const query = `select description, date_due, date_created, date_updated
		from task
		where id = $1`
	row := dbr.sdb.QueryRowContext(ctx, query, id)

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

func (dbr *DBRepo) Save(ctx context.Context, t Task) error {
	const query = `insert into "task" (id, description, date_due, date_created, date_updated)
		values ($1, $2, $3, $4, $5)`
	_, err := dbr.sdb.ExecContext(ctx,
		query,
		t.Id,
		t.Description,
		t.DateDue,
		t.DateCreated,
		t.DateUpdated)

	return err
}
