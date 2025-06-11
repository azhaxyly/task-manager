package in

import (
	"context"
	"time"

	"task-manager/internal/domain"
)

type GetTaskQuery struct {
	ID domain.TaskID
}

type TaskDTO struct {
	ID        domain.TaskID `json:"id"`
	Status    domain.Status `json:"status"`
	CreatedAt time.Time     `json:"created_at"`
	Duration  string        `json:"duration"`
	Result    *string       `json:"result,omitempty"`
	Error     *string       `json:"error,omitempty"`
}

type GetTaskUseCase interface {
	Handle(ctx context.Context, q GetTaskQuery) (TaskDTO, error)
}
