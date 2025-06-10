package http

import (
	"encoding/json"
	"net/http"
	"strings"
	"task-manager/internal/application/port/in"
	"task-manager/internal/common/logger"
	"task-manager/internal/domain"
	"time"
)

type TaskHandler struct {
	create in.CreateTaskUseCase
	get    in.GetTaskUseCase
	delete in.DeleteTaskUseCase
	list   in.ListTasksUseCase
}

func NewTaskHandler(
	create in.CreateTaskUseCase,
	get in.GetTaskUseCase,
	delete in.DeleteTaskUseCase,
	list in.ListTasksUseCase,
) *TaskHandler {
	return &TaskHandler{
		create,
		get,
		delete,
		list,
	}
}

func (h *TaskHandler) HandleTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	start := time.Now()
	logger.Info("Incoming %s %s", r.Method, r.URL.Path)

	switch r.Method {
	case http.MethodPost:
		h.handleCreate(w, r)
	case http.MethodGet:
		h.handleList(w, r)
	default:
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
	}

	elapsed := time.Since(start)
	logger.Info("Handled %s %s in %v", r.Method, r.URL.Path, elapsed)
}

func (h *TaskHandler) HandleTaskByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	start := time.Now()
	logger.Info("Incoming %s %s", r.Method, r.URL.Path)

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

	elapsed := time.Since(start)
	logger.Info("Handled %s %s in %v", r.Method, r.URL.Path, elapsed)
}

func (h *TaskHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	id, err := h.create.Handle(r.Context(), in.CreateTaskCommand{})
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":        id,
		"status":    domain.Pending,
		"createdAt": time.Now().UTC(),
	})
}

func (h *TaskHandler) handleGet(w http.ResponseWriter, r *http.Request, id domain.TaskID) {
	dto, err := h.get.Handle(r.Context(), in.GetTaskQuery{ID: id})
	if err != nil {
		if err == domain.ErrTaskNotFound {
			http.Error(w, `{"error":"task not found"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto)
}

func (h *TaskHandler) handleDelete(w http.ResponseWriter, r *http.Request, id domain.TaskID) {
	err := h.delete.Handle(r.Context(), in.DeleteTaskCommand{ID: id})
	if err != nil {
		if err == domain.ErrTaskNotFound {
			http.Error(w, `{"error":"task not found"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TaskHandler) handleList(w http.ResponseWriter, r *http.Request) {
	ids, err := h.list.Handle(r.Context(), in.ListTasksQuery{})
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ids)
}
