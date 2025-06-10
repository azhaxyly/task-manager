package service

import (
	"context"
	"task-manager/internal/application/port/in"
	"task-manager/internal/application/port/out"
	"task-manager/internal/common/logger"
	"task-manager/internal/domain"
)

type CreateTaskHandler struct {
	repo      out.TaskRepository
	scheduler out.TaskScheduler
	idGen     out.IDGenerator
}

func NewCreateTaskHandler(
	repo out.TaskRepository,
	scheduler out.TaskScheduler,
	idGen out.IDGenerator,
) *CreateTaskHandler {
	return &CreateTaskHandler{repo: repo, scheduler: scheduler, idGen: idGen}
}

func (h *CreateTaskHandler) Handle(ctx context.Context, cmd in.CreateTaskCommand) (domain.TaskID, error) {
	id := h.idGen.NewID()
	task := domain.NewTask(id)

	if err := h.repo.Save(ctx, task); err != nil {
		logger.Error("CreateTask failed: %v", err)
		return "", err
	}

	h.scheduler.Schedule(ctx, id)
	return id, nil
}
