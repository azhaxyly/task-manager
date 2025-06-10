package in

import (
	"context"
	"task-manager/internal/domain"
)

type DeleteTaskCommand struct {
	ID domain.TaskID
}

type DeleteTaskUseCase interface {
	Handle(ctx context.Context, cmd DeleteTaskCommand) error
}
