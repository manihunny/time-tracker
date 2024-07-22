package taskrepository

import (
	"context"
	"time"
)

type Task struct {
	ID                 int           `json:"id,omitempty" db:"id" example:"1"`
	UserID             int           `json:"user_id,omitempty" db:"user_id" example:"1"`
	Title              string        `json:"title,omitempty" db:"title" example:"Выполнить задачу 1"`
	StartedAt          time.Time     `json:"started_at,omitempty" db:"started_at" example:"2024-07-17T00:00:00Z"`
	FinishedAt         time.Time     `json:"finished_at,omitempty" db:"finished_at" example:"2024-07-17T00:00:00Z"`
	TimeSpentDuration  time.Duration `json:"time_spent,omitempty" swaggertype:"primitive,integer" example:"48393984418000"`
	TimeSpentFormatted string        `json:"time_spent_formatted,omitempty" example:"13h26m33.984418s"`
}

type TaskRepository interface {
	Create(ctx context.Context, t *Task) error
	Get(ctx context.Context, id string) (*Task, error)
	List(ctx context.Context) ([]*Task, error)
	Update(ctx context.Context, t *Task, id string) error
	PartialUpdate(ctx context.Context, a map[string]interface{}, id string) error
	Delete(ctx context.Context, id string) error
}
