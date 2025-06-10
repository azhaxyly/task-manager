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

type GetTaskHandler struct {
	repo out.TaskRepository
}

func NewGetTaskHandler(repo out.TaskRepository) *GetTaskHandler {
	return &GetTaskHandler{repo: repo}
}

func (h *GetTaskHandler) Handle(ctx context.Context, q in.GetTaskQuery) (in.TaskDTO, error) {
	t, err := h.repo.Find(ctx, q.ID)
	if err != nil {
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
		return err
	}

	switch t.Status {
	case domain.Pending, domain.Running:
		h.scheduler.Cancel(ctx, cmd.ID)
		if err2 := t.Cancel(); err2 != nil {
			return err2
		}
		return h.repo.Save(ctx, t)
	default:
		return h.repo.Delete(ctx, cmd.ID)
	}
}

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
