package taskrepository

import (
	"context"
	"time"
)

type Task struct {
	ID                 int           `json:"id,omitempty" db:"id"`
	UserID             int           `json:"user_id,omitempty" db:"user_id"`
	Title              string        `json:"title,omitempty" db:"title"`
	StartedAt          time.Time     `json:"started_at,omitempty" db:"started_at"`
	FinishedAt         time.Time     `json:"finished_at,omitempty" db:"finished_at"`
	TimeSpentDuration  time.Duration `json:"time_spent,omitempty"`
	TimeSpentFormatted string        `json:"time_spent_formatted,omitempty"`
}

type TaskRepository interface {
	Create(ctx context.Context, t *Task) error
	Get(ctx context.Context, id string) (*Task, error)
	List(ctx context.Context) ([]*Task, error)
	Update(ctx context.Context, t *Task, id string) error
	PartialUpdate(ctx context.Context, a map[string]interface{}, id string) error
	Delete(ctx context.Context, id string) error
}
