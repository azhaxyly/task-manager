package cmd

import (
	"log"
	"net/http"
	"os"
	myhttp "task-manager/internal/adapter/inbound/http"
	"task-manager/internal/adapter/outbound/idgen"
	"task-manager/internal/adapter/outbound/memstore"
	"task-manager/internal/application/service"
	"task-manager/internal/common/logger"

	_ "task-manager/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

func Run() {
	logger.Init(os.Stdout)

	repo := memstore.NewTaskRepository()
	scheduler := memstore.NewTaskScheduler(repo)

	uuidGen := idgen.NewUUIDGenerator()

	createH := service.NewCreateTaskHandler(repo, scheduler, uuidGen)
	getH := service.NewGetTaskHandler(repo)
	deleteH := service.NewDeleteTaskHandler(repo, scheduler)
	listH := service.NewListTasksHandler(repo)

	taskH := myhttp.NewTaskHandler(createH, getH, deleteH, listH)

	mux := myhttp.NewRouter(taskH)

	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	logger.Info("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
