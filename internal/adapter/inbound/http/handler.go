package http

import (
	"net/http"
	"task-manager/internal/application/port/in"
)

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

func (h *TaskHandler) handleTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
}
