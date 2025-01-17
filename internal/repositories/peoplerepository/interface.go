package peoplerepository

import (
	"context"
	"main/internal/repositories/taskrepository"
	"time"
)

type People struct {
	ID             int    `json:"-" db:"id" example:"1"`
	Surname        string `json:"surname,omitempty" db:"surname" qp:"surname" example:"Иванов"`
	Name           string `json:"name,omitempty" db:"name" qp:"name" example:"Иван"`
	Patronymic     string `json:"patronymic,omitempty" db:"patronymic" qp:"patronymic" example:"Иванович"`
	Address        string `json:"address,omitempty" db:"address" qp:"address" example:"г. Москва, ул. Ленина, д. 5, кв. 1"`
	PassportNumber string `json:"passport_number,omitempty" db:"passport_number" qp:"passport_number" example:"1234 567890"`
}

type PeopleRepository interface {
	List(ctx context.Context, limit int, offset int, filter *People) ([]*People, error)
	Get(ctx context.Context, id string) (*People, error)
	Create(ctx context.Context, p *People) error
	Update(ctx context.Context, p *People, id string) error
	PartialUpdate(ctx context.Context, a map[string]interface{}, id string) error
	Delete(ctx context.Context, id string) error
	StartNewTaskForUser(ctx context.Context, userID int, title string) error
	FinishAllUserTasks(ctx context.Context, userID int) error
	TaskStatistics(ctx context.Context, date_from time.Time, date_to time.Time) ([]*taskrepository.Task, error)
}
