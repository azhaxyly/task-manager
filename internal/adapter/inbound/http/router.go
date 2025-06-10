package http

import "net/http"

// RegisterRoutes регистрирует эндпоинты в mux.
func (h *TaskHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/tasks", h.handleTasks) // POST /tasks - создать таску
	// GET /tasks/{id} - получить таску по айди
	// DELETE /tasks/{id} - удалить таску по айди
}

func NewRouter(handler *TaskHandler) *http.ServeMux {
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	return mux
}
