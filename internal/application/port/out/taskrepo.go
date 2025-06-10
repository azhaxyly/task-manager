package out

import (
	"context"
	"task-manager/internal/domain"
)

type TaskRepository interface {
	Add(ctx context.Context, task *domain.Task) error
	Save(ctx context.Context, t *domain.Task) error
	Find(ctx context.Context, id domain.TaskID) (*domain.Task, error)
	Delete(ctx context.Context, id domain.TaskID) error
}
