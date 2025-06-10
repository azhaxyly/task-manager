package service

import (
	"context"
	"strconv"
	"task-manager/internal/application/port/in"
	"task-manager/internal/application/port/out"
	"task-manager/internal/domain"
	"time"
)

type CreateTaskHandler struct {
	repo      out.TaskRepository
	scheduler out.TaskScheduler
}

func NewCreateTaskHandler(repo out.TaskRepository, scheduler out.TaskScheduler) *CreateTaskHandler {
	return &CreateTaskHandler{repo: repo, scheduler: scheduler}
}

func (h *CreateTaskHandler) Handle(ctx context.Context, cmd in.CreateTaskCommand) (domain.TaskID, error) {
	id := domain.TaskID(strconv.FormatInt(time.Now().UTC().UnixNano(), 10))
	task := domain.NewTask(id)

	if err := h.repo.Save(ctx, task); err != nil {
		return "", err
	}

	h.scheduler.Schedule(ctx, id)
	return id, nil
}
