package cmd

import (
	"log"
	"net/http"
	myhttp "task-manager/internal/adapter/inbound/http"
	"task-manager/internal/adapter/outbound/memstore"
	"task-manager/internal/application/service"
)

func Run() {
	repo := memstore.NewTaskRepository()
	scheduler := memstore.NewTaskScheduler(repo)

	createHandler := service.NewCreateTaskHandler(repo, scheduler)
	getHandler := service.NewGetTaskHandler(repo)
	deleteHandler := service.NewDeleteTaskHandler(repo, scheduler)

	taskHandler := myhttp.NewTaskHandler(createHandler, getHandler, deleteHandler)

	mux := myhttp.NewRouter(taskHandler)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
