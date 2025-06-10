package http

import (
	"encoding/json"
	"net/http"
	"task-manager/internal/application/port/in"
	"task-manager/internal/domain"
)

type TaskHandler struct {
	create in.CreateTaskUseCase
	// get    in.GetTaskUseCase
	// delete in.DeleteTaskUseCase
}

func NewTaskHandler(
	create in.CreateTaskUseCase,
	// get in.GetTaskUseCase,
	// delete in.DeleteTaskUseCase,
) *TaskHandler {
	return &TaskHandler{
		create: create,
		// get: get,
		// delete: delete,
	}
}

func (h *TaskHandler) handleTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	id, err := h.create.Handle(r.Context(), in.CreateTaskCommand{})
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":     id,
		"status": domain.Pending,
	})
}
