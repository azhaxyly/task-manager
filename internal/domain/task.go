package domain

import "time"

type TaskID string

type Task struct {
	ID         TaskID
	CreatedAt  time.Time
	StartedAt  *time.Time
	FinishedAt *time.Time
	Status     Status
	Result     *string
	Err        *string
}

func NewTask(id TaskID) *Task {
	return &Task{
		ID:        id,
		CreatedAt: time.Now().UTC(),
		Status:    Pending,
	}
}
