package domain

import (
	"errors"
	"time"
)

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

func (t *Task) Start() error {
	if t.Status != Pending {
		return errors.New("task can't start, status isn't 'pending'")
	}
	now := time.Now().UTC()
	t.StartedAt = &now
	t.Status = Running
	return nil
}
