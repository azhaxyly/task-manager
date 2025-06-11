package service

import (
	"context"
	"fmt"
	"time"

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

	raw := t.Duration().Truncate(time.Second)
	mins := int(raw / time.Minute)
	secs := int((raw % time.Minute) / time.Second)
	formatted := fmt.Sprintf("%dm%ds", mins, secs)

	return in.TaskDTO{
		ID:        t.ID,
		Status:    t.Status,
		CreatedAt: t.CreatedAt,
		Duration:  formatted,
		Result:    t.Result,
		Error:     t.Err,
	}, nil
}
