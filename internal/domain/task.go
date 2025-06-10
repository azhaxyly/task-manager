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

func (t *Task) Complete(result string) error {
	if t.Status != Running {
		return errors.New("task can't complete, status isn't 'running'")
	}
	now := time.Now().UTC()
	t.FinishedAt = &now
	t.Result = &result
	t.Status = Success
	return nil
}

// пока не используется, на будущее
func (t *Task) Fail(err error) error {
	if t.Status != Running {
		return errors.New("task can't fail, status isn't 'running'")
	}
	now := time.Now().UTC()
	t.FinishedAt = &now
	msg := err.Error()
	t.Err = &msg
	t.Status = Failed
	return nil
}

func (t *Task) Cancel() error {
	if t.Status != Pending && t.Status != Running {
		return errors.New("task can't cancel, status isn't 'pending' or 'running'")
	}
	now := time.Now().UTC()
	t.FinishedAt = &now
	msg := "canceled"
	t.Err = &msg
	t.Status = Canceled
	return nil
}

// в наносеках
func (t *Task) Duration() time.Duration {
	if t.StartedAt == nil {
		return 0
	}
	if t.FinishedAt == nil {
		return time.Since(*t.StartedAt)
	}
	return t.FinishedAt.Sub(*t.StartedAt)
}
