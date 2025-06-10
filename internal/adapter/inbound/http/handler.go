package http

import (
	"encoding/json"
	"net/http"
	"strings"
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

func (h *TaskHandler) handleTaskByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := strings.TrimPrefix(r.URL.Path, "/tasks/")
	if id == "" {
		http.Error(w, `{"error":"missing task id"}`, http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGet(w, r, domain.TaskID(id))
	case http.MethodDelete:
		h.handleDelete(w, r, domain.TaskID(id))
	default:
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func (h *TaskHandler) handleGet(w http.ResponseWriter, r *http.Request, id domain.TaskID) {}

func (h *TaskHandler) handleDelete(w http.ResponseWriter, r *http.Request, id domain.TaskID) {}
