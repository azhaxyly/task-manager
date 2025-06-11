package service

import (
	"context"

	"task-manager/internal/application/port/in"
	"task-manager/internal/application/port/out"
	"task-manager/internal/common/logger"
	"task-manager/internal/domain"
)

type DeleteTaskHandler struct {
	repo      out.TaskRepository
	scheduler out.TaskScheduler
}

func NewDeleteTaskHandler(repo out.TaskRepository, scheduler out.TaskScheduler) *DeleteTaskHandler {
	return &DeleteTaskHandler{repo: repo, scheduler: scheduler}
}

func (h *DeleteTaskHandler) Handle(ctx context.Context, cmd in.DeleteTaskCommand) error {
	t, err := h.repo.Find(ctx, cmd.ID)
	if err != nil {
		logger.Error("DeleteTask failed: %v", err)
		return err
	}

	switch t.Status {
	case domain.Pending, domain.Running:
		h.scheduler.Cancel(ctx, cmd.ID)
		if err2 := t.Cancel(); err2 != nil {
			logger.Error("Failed to cancel task %s: %v", cmd.ID, err2)
			return err2
		}
		return h.repo.Save(ctx, t)
	default:
		return h.repo.Delete(ctx, cmd.ID)
	}
}
