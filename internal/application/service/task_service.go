package service

import "task-manager/internal/application/port/out"

type CreateTaskHandler struct {
	repo      out.TaskRepository
	scheduler out.TaskScheduler
}

func NewCreateTaskHandler(repo out.TaskRepository, scheduler out.TaskScheduler) *CreateTaskHandler {
	return &CreateTaskHandler{repo: repo, scheduler: scheduler}
}
