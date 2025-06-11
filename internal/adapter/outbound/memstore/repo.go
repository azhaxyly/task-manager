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

func (r *TaskRepository) Find(ctx context.Context, id domain.TaskID) (*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	original, ok := r.data[id]
	if !ok {
		return nil, domain.ErrTaskNotFound
	}
	return cloneTask(original), nil
}

func (r *TaskRepository) Delete(ctx context.Context, id domain.TaskID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[id]; !ok {
		return domain.ErrTaskNotFound
	}
	delete(r.data, id)
	return nil
}

func (r *TaskRepository) List(ctx context.Context) ([]domain.TaskID, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ids := make([]domain.TaskID, 0, len(r.data))
	for id := range r.data {
		ids = append(ids, id)
	}
	return ids, nil
}

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
