package service

import (
	"context"
	"task-manager/internal/application/port/in"
	"task-manager/internal/application/port/out"
	"task-manager/internal/common/logger"
)

type GetTaskHandler struct {
	repo out.TaskRepository
}

func NewGetTaskHandler(repo out.TaskRepository) *GetTaskHandler {
	return &GetTaskHandler{repo: repo}
}

func (h *GetTaskHandler) Handle(ctx context.Context, q in.GetTaskQuery) (in.TaskDTO, error) {
	t, err := h.repo.Find(ctx, q.ID)
	if err != nil {
		logger.Error("GetTask failed: %v", err)
		return in.TaskDTO{}, err
	}
	return in.TaskDTO{
		ID:        t.ID,
		Status:    t.Status,
		CreatedAt: t.CreatedAt,
		Duration:  t.Duration(),
		Result:    t.Result,
		Error:     t.Err,
	}, nil
}
