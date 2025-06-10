package in

import (
	"context"
	"task-manager/internal/domain"
	"time"
)

type GetTaskQuery struct {
	ID domain.TaskID
}

type TaskDTO struct {
	ID        domain.TaskID `json:"id"`
	Status    domain.Status `json:"status"`
	CreatedAt time.Time     `json:"created_at"`
	Duration  time.Duration `json:"duration"`
	Result    *string       `json:"result,omitempty"`
	Error     *string       `json:"error,omitempty"`
}

type GetTaskUseCase interface {
	Handle(ctx context.Context, q GetTaskQuery) (TaskDTO, error)
}
