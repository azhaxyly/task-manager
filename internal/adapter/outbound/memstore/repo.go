package memstore

import (
	"context"
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

func (r *TaskRepository) Find(ctx context.Context, id domain.TaskID) (*domain.Task, error)

func (r *TaskRepository) Delete(ctx context.Context, id domain.TaskID) error
