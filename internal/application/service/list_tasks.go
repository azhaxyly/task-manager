package service

import (
	"context"
	"task-manager/internal/application/port/in"
	"task-manager/internal/application/port/out"
)

type ListTasksHandler struct {
	repo out.TaskRepository
}

func NewListTasksHandler(repo out.TaskRepository) *ListTasksHandler {
	return &ListTasksHandler{repo: repo}
}

func (h *ListTasksHandler) Handle(ctx context.Context, _ in.ListTasksQuery) ([]in.TaskSummaryDTO, error) {
	ids, _ := h.repo.List(ctx)
	summaries := make([]in.TaskSummaryDTO, len(ids))
	for i, id := range ids {
		t, _ := h.repo.Find(ctx, id)
		summaries[i] = in.TaskSummaryDTO{
			ID:        t.ID,
			Status:    t.Status,
			CreatedAt: t.CreatedAt,
		}
	}
	return summaries, nil
}
