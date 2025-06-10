package memstore

import (
	"context"
	"errors"
	"sync"
	"task-manager/internal/domain"
)

type TaskRepository struct {
	mu   sync.RWMutex
	data map[domain.TaskID]*domain.Task
}

func NewTaskRepository() *TaskRepository {
	return &TaskRepository{
		data: make(map[domain.TaskID]*domain.Task),
	}
}

func (r *TaskRepository) Save(ctx context.Context, t *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[t.ID] = t
	return nil
}

func (r *TaskRepository) Find(ctx context.Context, id domain.TaskID) (*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	original, ok := r.data[id]
	if !ok {
		return nil, errors.New("task not found")
	}
	// TODO: всеравно нужно делать копию, чтобы не было дата рейса
	return cloneTask(original), nil
}

func (r *TaskRepository) Delete(ctx context.Context, id domain.TaskID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[id]; !ok {
		return errors.New("task not found")
	}
	delete(r.data, id)
	return nil
}

// TODO: нужно доделать
func cloneTask(t *domain.Task) *domain.Task {
	c := *t

	if t.StartedAt != nil {
		temp := *t.StartedAt
		c.StartedAt = &temp
	}
	if t.FinishedAt != nil {
		temp := *t.FinishedAt
		c.FinishedAt = &temp
	}
	if t.Result != nil {
		temp := *t.Result
		c.Result = &temp
	}
	if t.Err != nil {
		temp := *t.Err
		c.Err = &temp
	}

	return &c
}
