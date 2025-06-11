package in

import (
	"context"
	"time"

	"task-manager/internal/domain"
)

type ListTasksQuery struct{}

type TaskSummaryDTO struct {
	ID        domain.TaskID `json:"id"`
	Status    domain.Status `json:"status"`
	CreatedAt time.Time     `json:"created_at"`
	Duration  string        `json:"duration"`
}
type ListTasksUseCase interface {
	Handle(ctx context.Context, q ListTasksQuery) ([]TaskSummaryDTO, error)
}
