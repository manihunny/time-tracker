package tasksqlx

import (
	"context"
	"fmt"
	"main/internal/repositories/taskrepository"
	"strings"

	"github.com/jmoiron/sqlx"
)

type TaskSqlx struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *TaskSqlx {
	return &TaskSqlx{
		db: db,
	}
}

func (r *TaskSqlx) List(ctx context.Context) ([]*taskrepository.Task, error) {
	const q = `
		select * from tasks
	`
	list := []*taskrepository.Task{}
	err := r.db.SelectContext(ctx, &list, q)
	return list, err
}

func (r *TaskSqlx) Get(ctx context.Context, id string) (*taskrepository.Task, error) {
	const q = `
		select * from tasks where id = $1
	`
	t := new(taskrepository.Task)
	err := r.db.GetContext(ctx, t, q, id)
	return t, err
}

func (r *TaskSqlx) Create(ctx context.Context, t *taskrepository.Task) error {
	const q = `
		insert into tasks (id, user_id, title, started_at, finished_at) 
			values (:id, :user_id, :title, :started_at, :finished_at)
	`
	_, err := r.db.NamedExecContext(ctx, q, t)
	return err
}

func (r *TaskSqlx) Update(ctx context.Context, t *taskrepository.Task, id string) error {
	const q = `
		update tasks set id = $1, user_id = $2, title = $3, started_at = $4, finished_at = $5 
			where id = $6
	`
	_, err := r.db.ExecContext(ctx, q, t.ID, t.UserID, t.Title, t.StartedAt, t.FinishedAt, id)
	return err
}

func (r *TaskSqlx) PartialUpdate(ctx context.Context, taskData map[string]interface{}, id string) error {
	q := `
		update tasks set 
	`
	for key, value := range taskData {
		q += fmt.Sprintf(`%v = '%v',`, key, value)
	}
	q = strings.TrimRight(q, ",") + fmt.Sprintf(` WHERE id = '%v'`, id)
	_, err := r.db.ExecContext(ctx, q)
	return err
}

func (r *TaskSqlx) Delete(ctx context.Context, id string) error {
	const q = `
		delete from tasks where id = $1
	`
	_, err := r.db.ExecContext(ctx, q, id)
	return err
}
