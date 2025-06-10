package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	myhttp "task-manager/internal/adapter/inbound/http"
	"task-manager/internal/application/port/in"
	"task-manager/internal/domain"
)

type stubCreate struct{ id domain.TaskID }

func (s stubCreate) Handle(ctx context.Context, cmd in.CreateTaskCommand) (domain.TaskID, error) {
	return s.id, nil
}

type stubGet struct {
	dto in.TaskDTO
	err error
}

func (s stubGet) Handle(ctx context.Context, q in.GetTaskQuery) (in.TaskDTO, error) {
	return s.dto, s.err
}

type stubDelete struct{ err error }

func (s stubDelete) Handle(ctx context.Context, cmd in.DeleteTaskCommand) error {
	return s.err
}

type stubList struct {
	list []in.TaskSummaryDTO
	err  error
}

func (s stubList) Handle(ctx context.Context, q in.ListTasksQuery) ([]in.TaskSummaryDTO, error) {
	return s.list, s.err
}

func TestTaskHandler_Create(t *testing.T) {
	h := myhttp.NewTaskHandler(
		stubCreate{"ID123"},
		stubGet{},
		stubDelete{},
		stubList{},
	)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/tasks", nil)

	h.HandleTasks(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("Create: code=%d; want %d", rec.Code, http.StatusCreated)
	}
	var body map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatal("Create: decode body:", err)
	}
	if body["id"] != "ID123" {
		t.Errorf("Create: id=%v; want ID123", body["id"])
	}
	if body["status"] != string(domain.Pending) {
		t.Errorf("Create: status=%v; want %v", body["status"], domain.Pending)
	}
}

func TestTaskHandler_List(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	dto := in.TaskSummaryDTO{
		ID:        "X",
		Status:    domain.Pending,
		CreatedAt: now,
		Duration:  "0m0s",
	}
	h := myhttp.NewTaskHandler(
		stubCreate{},
		stubGet{},
		stubDelete{},
		stubList{[]in.TaskSummaryDTO{dto}, nil},
	)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)

	h.HandleTasks(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("List: code=%d; want %d", rec.Code, http.StatusOK)
	}
	var got []in.TaskSummaryDTO
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatal("List: decode body:", err)
	}
	if len(got) != 1 {
		t.Fatalf("List: len=%d; want 1", len(got))
	}
	if got[0].ID != "X" || got[0].Status != domain.Pending || !got[0].CreatedAt.Equal(now) || got[0].Duration != "0m0s" {
		t.Errorf("List: got %+v; want %+v", got[0], dto)
	}
}

func TestTaskHandler_GetSuccess(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	expected := in.TaskDTO{
		ID:        "ID1",
		Status:    domain.Success,
		CreatedAt: now,
		Duration:  "1m2s",
		Result:    ptr("done"),
		Error:     nil,
	}
	h := myhttp.NewTaskHandler(
		stubCreate{},
		stubGet{expected, nil},
		stubDelete{},
		stubList{},
	)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/tasks/ID1", nil)

	h.HandleTaskByID(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("GetSuccess: code=%d; want %d", rec.Code, http.StatusOK)
	}
	var dto in.TaskDTO
	if err := json.NewDecoder(rec.Body).Decode(&dto); err != nil {
		t.Fatal("GetSuccess: decode body:", err)
	}
	if dto.ID != expected.ID ||
		dto.Status != expected.Status ||
		!dto.CreatedAt.Equal(expected.CreatedAt) ||
		dto.Duration != expected.Duration ||
		(dto.Result == nil) != (expected.Result == nil) ||
		(dto.Result != nil && *dto.Result != *expected.Result) {
		t.Errorf("GetSuccess: got %+v; want %+v", dto, expected)
	}
}

func TestTaskHandler_GetNotFound(t *testing.T) {
	h := myhttp.NewTaskHandler(
		stubCreate{},
		stubGet{in.TaskDTO{}, domain.ErrTaskNotFound},
		stubDelete{},
		stubList{},
	)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/tasks/XXX", nil)

	h.HandleTaskByID(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("GetNotFound: code=%d; want %d", rec.Code, http.StatusNotFound)
	}
}

func TestTaskHandler_DeleteSuccess(t *testing.T) {
	h := myhttp.NewTaskHandler(
		stubCreate{},
		stubGet{},
		stubDelete{nil},
		stubList{},
	)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/tasks/ID1", nil)

	h.HandleTaskByID(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("DeleteSuccess: code=%d; want %d", rec.Code, http.StatusNoContent)
	}
}

func TestTaskHandler_DeleteNotFound(t *testing.T) {
	h := myhttp.NewTaskHandler(
		stubCreate{},
		stubGet{},
		stubDelete{domain.ErrTaskNotFound},
		stubList{},
	)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/tasks/ID1", nil)

	h.HandleTaskByID(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("DeleteNotFound: code=%d; want %d", rec.Code, http.StatusNotFound)
	}
}

func TestTaskHandler_MethodNotAllowed_HandleTasks(t *testing.T) {
	h := myhttp.NewTaskHandler(stubCreate{}, stubGet{}, stubDelete{}, stubList{})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/tasks", bytes.NewBuffer(nil))

	h.HandleTasks(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("MethodNotAllowed POST/GET: code=%d; want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}

func TestTaskHandler_MethodNotAllowed_HandleTaskByID(t *testing.T) {
	h := myhttp.NewTaskHandler(stubCreate{}, stubGet{}, stubDelete{}, stubList{})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, "/tasks/ID1", nil)

	h.HandleTaskByID(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("MethodNotAllowed GET/DELETE: code=%d; want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}

func ptr(s string) *string {
	return &s
}
