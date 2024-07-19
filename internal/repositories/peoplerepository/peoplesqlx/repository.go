package peoplesqlx

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"main/internal/repositories/peoplerepository"
	"main/internal/repositories/taskrepository"
	"strings"
	"time"
)

type PeopleSqlx struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *PeopleSqlx {
	return &PeopleSqlx{
		db: db,
	}
}

func (r *PeopleSqlx) List(ctx context.Context, limit int, offset int, filterFields *peoplerepository.People) ([]*peoplerepository.People, error) {
	q := `
		select * from people 
			where lower(surname) like lower($1) 
			and lower(name) like lower($2) 
			and lower(patronymic) like lower($3) 
			and lower(address) like lower($4) 
			and lower(passport_number) like lower($5) 
			order by id
			limit $6
			offset $7
	`
	list := []*peoplerepository.People{}
	err := r.db.SelectContext(ctx, &list, q, "%"+filterFields.Surname+"%", "%"+filterFields.Name+"%", "%"+filterFields.Patronymic+"%", "%"+filterFields.Address+"%", "%"+filterFields.PassportNumber+"%", limit, offset)
	return list, err
}

func (r *PeopleSqlx) Get(ctx context.Context, id string) (*peoplerepository.People, error) {
	const q = `
		select * from people where id = $1
	`
	p := new(peoplerepository.People)
	err := r.db.GetContext(ctx, p, q, id)
	return p, err
}

func (r *PeopleSqlx) Create(ctx context.Context, p *peoplerepository.People) error {
	const q = `
		insert into people (surname, name, patronymic, address, passport_number) 
			values (:surname, :name, :patronymic, :address, :passport_number)
	`
	_, err := r.db.NamedExecContext(ctx, q, p)
	return err
}

func (r *PeopleSqlx) Update(ctx context.Context, p *peoplerepository.People, id string) error {
	const q = `
		update people set surname = $1, name = $2, patronymic = $3, address = $4, passport_number = $5 
			where id = $6
	`
	_, err := r.db.ExecContext(ctx, q, p.Surname, p.Name, p.Patronymic, p.Address, p.PassportNumber, id)
	return err
}

func (r *PeopleSqlx) PartialUpdate(ctx context.Context, PeopleData map[string]interface{}, id string) error {
	q := `
		update people set 
	`
	for key, value := range PeopleData {
		q += fmt.Sprintf(`%v = '%v',`, key, value)
	}
	q = strings.TrimRight(q, ",") + fmt.Sprintf(` WHERE id = '%v'`, id)
	_, err := r.db.ExecContext(ctx, q)
	return err
}

func (r *PeopleSqlx) Delete(ctx context.Context, id string) error {
	const q = `
		delete from people where id = $1
	`
	_, err := r.db.ExecContext(ctx, q, id)
	return err
}

func (r *PeopleSqlx) StartNewTaskForUser(ctx context.Context, userID int, title string) error {
	const q = `
		insert into tasks (user_id, title, started_at, finished_at) 
			values ($1, $2, now(), null)
	`
	_, err := r.db.ExecContext(ctx, q, userID, title)
	return err
}

func (r *PeopleSqlx) FinishAllUserTasks(ctx context.Context, userID int) error {
	const q = `
		update tasks set finished_at = now()
			where user_id = $1 and finished_at is null
	`
	_, err := r.db.ExecContext(ctx, q, userID)
	return err
}

func (r *PeopleSqlx) TaskStatistics(ctx context.Context, date_from time.Time, date_to time.Time) ([]*taskrepository.Task, error) {
	const q = `
		select * from tasks
			where (started_at between $1 and $2) or (finished_at between $1 and $2)
	`
	list := []*taskrepository.Task{}
	err := r.db.SelectContext(ctx, &list, q, date_from, date_to)
	return list, err
}
