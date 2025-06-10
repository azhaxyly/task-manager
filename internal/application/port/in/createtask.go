package in

import (
	"context"
	"task-manager/internal/domain"
)

// TODO: пока хз
type CreateTaskCommand struct{}

type CreateTaskUseCase interface {
	Handle(ctx context.Context, cmd CreateTaskCommand) (domain.TaskID, error)
}
