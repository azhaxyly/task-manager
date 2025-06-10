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
	return original, nil
}

func (r *TaskRepository) Delete(ctx context.Context, id domain.TaskID) error
