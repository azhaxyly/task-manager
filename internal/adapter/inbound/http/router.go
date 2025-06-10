package http

import "net/http"

// RegisterRoutes регистрирует эндпоинты в mux.
func (h *TaskHandler) RegisterRoutes(mux *http.ServeMux) {
	// POST /tasks - создать таск
	// GET /tasks/{id} - получить таску по айди
	// DELETE /tasks/{id} - удалить таску по айди
}
