package http

import "task-manager/internal/application/port/in"

type TaskHandler struct {
	create in.CreateTaskUseCase
	get    in.GetTaskUseCase
	delete in.DeleteTaskUseCase
}

func NewTaskHandler(
	create in.CreateTaskUseCase,
	get in.GetTaskUseCase,
	delete in.DeleteTaskUseCase,
) *TaskHandler {
	return &TaskHandler{create: create, get: get, delete: delete}
}
