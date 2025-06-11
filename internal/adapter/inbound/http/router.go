package http

import (
	"net/http"
)

func (h *TaskHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/tasks", h.HandleTasks)
	// POST /tasks - создать таску  ***РАБОТАЕТ***
	// GET /tasks - получить список тасок ***РАБОТАЕТ***
	mux.HandleFunc("/tasks/", h.HandleTaskByID)
	// GET /tasks/{id} - получить таску по айди ***РАБОТАЕТ***
	// DELETE /tasks/{id} - удалить таску по айди ***РАБОТАЕТ***
}

func NewRouter(handler *TaskHandler) *http.ServeMux {
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	return mux
}
