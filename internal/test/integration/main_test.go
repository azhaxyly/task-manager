package main_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	httpadp "task-manager/internal/adapter/inbound/http"
	idgen "task-manager/internal/adapter/outbound/idgen"
	memstore "task-manager/internal/adapter/outbound/memstore"
	"task-manager/internal/application/service"
	"task-manager/internal/common/logger"
	"task-manager/internal/domain"
)

type ImmediateScheduler struct {
	repo *memstore.TaskRepository
}

func NewImmediateScheduler(repo *memstore.TaskRepository) *ImmediateScheduler {
	return &ImmediateScheduler{repo}
}
func (s *ImmediateScheduler) Schedule(ctx context.Context, id domain.TaskID) {
	t, _ := s.repo.Find(ctx, id)
	_ = t.Start()
	_ = t.Complete("ok")
	s.repo.Save(ctx, t)
}
func (s *ImmediateScheduler) Cancel(ctx context.Context, id domain.TaskID) {
	t, err := s.repo.Find(ctx, id)
	if err != nil {
		return
	}
	_ = t.Cancel()
	s.repo.Save(ctx, t)
}

func TestIntegration_TaskLifecycle(t *testing.T) {
	logger.Init(io.Discard)

	repo := memstore.NewTaskRepository()
	scheduler := NewImmediateScheduler(repo)
	uuidGen := idgen.NewUUIDGenerator()

	createH := service.NewCreateTaskHandler(repo, scheduler, uuidGen)
	getH := service.NewGetTaskHandler(repo)
	deleteH := service.NewDeleteTaskHandler(repo, scheduler)
	listH := service.NewListTasksHandler(repo)

	handler := httpadp.NewTaskHandler(createH, getH, deleteH, listH)
	server := httptest.NewServer(httpadp.NewRouter(handler))
	defer server.Close()

	resp, err := http.Post(server.URL+"/tasks", "application/json", nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("POST code=%d; want %d", resp.StatusCode, http.StatusCreated)
	}
	var createBody struct{ ID string }
	json.NewDecoder(resp.Body).Decode(&createBody)
	if createBody.ID == "" {
		t.Fatal("empty ID")
	}

	resp2, err := http.Get(server.URL + "/tasks/" + createBody.ID)
	if err != nil {
		t.Fatal(err)
	}
	var dto struct {
		Status   domain.Status `json:"status"`
		Duration string        `json:"duration"`
	}
	json.NewDecoder(resp2.Body).Decode(&dto)
	if dto.Status != domain.Success {
		t.Errorf("status=%s; want %s", dto.Status, domain.Success)
	}
	if dto.Duration == "" {
		t.Error("empty duration")
	}

	req, _ := http.NewRequest(http.MethodDelete, server.URL+"/tasks/"+createBody.ID, nil)
	resp3, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp3.StatusCode != http.StatusNoContent {
		t.Errorf("DELETE code=%d; want %d", resp3.StatusCode, http.StatusNoContent)
	}

	resp4, err := http.Get(server.URL + "/tasks/" + createBody.ID)
	if err != nil {
		t.Fatal(err)
	}
	if resp4.StatusCode != http.StatusNotFound {
		t.Errorf("after delete GET code=%d; want %d", resp4.StatusCode, http.StatusNotFound)
	}
}
