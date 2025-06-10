package service

import (
	"context"
	"fmt"
	"task-manager/internal/application/port/in"
	"task-manager/internal/application/port/out"
	"time"
)

type ListTasksHandler struct {
	repo out.TaskRepository
}

func NewListTasksHandler(repo out.TaskRepository) *ListTasksHandler {
	return &ListTasksHandler{repo: repo}
}

func (h *ListTasksHandler) Handle(ctx context.Context, _ in.ListTasksQuery) ([]in.TaskSummaryDTO, error) {
	ids, err := h.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	summaries := make([]in.TaskSummaryDTO, len(ids))
	for i, id := range ids {
		t, err := h.repo.Find(ctx, id)
		if err != nil {
			return nil, err
		}

		raw := t.Duration().Truncate(time.Second)
		mins := int(raw / time.Minute)
		secs := int((raw % time.Minute) / time.Second)
		formatted := fmt.Sprintf("%dm%ds", mins, secs)

		summaries[i] = in.TaskSummaryDTO{
			ID:        t.ID,
			Status:    t.Status,
			CreatedAt: t.CreatedAt,
			Duration:  formatted,
		}
	}
	return summaries, nil
}
