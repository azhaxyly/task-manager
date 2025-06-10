package service_test

import (
	"context"
	"testing"

	"task-manager/internal/application/port/in"
	"task-manager/internal/application/port/out"
	"task-manager/internal/application/service"
	"task-manager/internal/domain"
)

type fakeRepo struct {
	saved   map[domain.TaskID]*domain.Task
	deleted []domain.TaskID
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{saved: make(map[domain.TaskID]*domain.Task)}
}

func (r *fakeRepo) Save(ctx context.Context, t *domain.Task) error {
	r.saved[t.ID] = t
	return nil
}
func (r *fakeRepo) Find(ctx context.Context, id domain.TaskID) (*domain.Task, error) {
	t, ok := r.saved[id]
	if !ok {
		return nil, domain.ErrTaskNotFound
	}
	return t, nil
}
func (r *fakeRepo) Delete(ctx context.Context, id domain.TaskID) error {
	if _, ok := r.saved[id]; !ok {
		return domain.ErrTaskNotFound
	}
	delete(r.saved, id)
	r.deleted = append(r.deleted, id)
	return nil
}
func (r *fakeRepo) List(ctx context.Context) ([]domain.TaskID, error) {
	ids := make([]domain.TaskID, 0, len(r.saved))
	for id := range r.saved {
		ids = append(ids, id)
	}
	return ids, nil
}

type fakeScheduler struct {
	scheduled []domain.TaskID
	canceled  []domain.TaskID
}

func newFakeScheduler() *fakeScheduler {
	return &fakeScheduler{}
}
func (s *fakeScheduler) Schedule(ctx context.Context, id domain.TaskID) {
	s.scheduled = append(s.scheduled, id)
}
func (s *fakeScheduler) Cancel(ctx context.Context, id domain.TaskID) {
	s.canceled = append(s.canceled, id)
}

func TestCreateTaskHandler_Handle(t *testing.T) {
	repo := newFakeRepo()
	sched := newFakeScheduler()
	idGen := &struct{ out.IDGenerator }{IDGenerator: out.NewStaticID("X")}
	h := service.NewCreateTaskHandler(repo, sched, idGen)

	id, err := h.Handle(context.Background(), in.CreateTaskCommand{})
	if err != nil {
		t.Fatalf("Handle error: %v", err)
	}
	if id != domain.TaskID("X") {
		t.Errorf("id=%s; want X", id)
	}
	if _, ok := repo.saved[id]; !ok {
		t.Errorf("task not saved in repo")
	}
	if len(sched.scheduled) != 1 || sched.scheduled[0] != id {
		t.Errorf("scheduler.Schedule not called with %s", id)
	}
}

func TestGetTaskHandler_Handle(t *testing.T) {
	repo := newFakeRepo()
	id := domain.TaskID("id1")
	task := domain.NewTask(id)
	task.Start()
	task.Complete("ok")
	repo.Save(context.Background(), task)

	h := service.NewGetTaskHandler(repo)
	dto, err := h.Handle(context.Background(), in.GetTaskQuery{ID: id})
	if err != nil {
		t.Fatalf("Handle error: %v", err)
	}
	if dto.ID != id || dto.Status != domain.Success {
		t.Errorf("got %+v; want ID=%s Status=%s", dto, id, domain.Success)
	}
	if dto.Duration == "0" {
		t.Errorf("Duration zero; want >0")
	}
}

func TestDeleteTaskHandler_Handle(t *testing.T) {
	repo := newFakeRepo()
	sched := newFakeScheduler()
	handler := service.NewDeleteTaskHandler(repo, sched)

	id1 := domain.TaskID("p")
	t1 := domain.NewTask(id1)
	repo.Save(context.Background(), t1)
	err := handler.Handle(context.Background(), in.DeleteTaskCommand{ID: id1})
	if err != nil {
		t.Fatalf("Delete pending error: %v", err)
	}
	if len(sched.canceled) != 1 || sched.canceled[0] != id1 {
		t.Errorf("Cancel not called for %s", id1)
	}
	t1Saved, _ := repo.Find(context.Background(), id1)
	if t1Saved.Status != domain.Canceled {
		t.Errorf("Status=%s; want %s", t1Saved.Status, domain.Canceled)
	}

	id2 := domain.TaskID("c")
	t2 := domain.NewTask(id2)
	t2.Start()
	t2.Complete("ok")
	repo.Save(context.Background(), t2)
	err = handler.Handle(context.Background(), in.DeleteTaskCommand{ID: id2})
	if err != nil {
		t.Fatalf("Delete completed error: %v", err)
	}
	if _, err2 := repo.Find(context.Background(), id2); err2 != domain.ErrTaskNotFound {
		t.Errorf("Completed task not deleted; err=%v", err2)
	}
}

func TestListTasksHandler_Handle(t *testing.T) {
	repo := newFakeRepo()
	handler := service.NewListTasksHandler(repo)

	idA, idB := domain.TaskID("A"), domain.TaskID("B")
	repo.Save(context.Background(), domain.NewTask(idA))
	repo.Save(context.Background(), domain.NewTask(idB))

	summaries, err := handler.Handle(context.Background(), in.ListTasksQuery{})
	if err != nil {
		t.Fatalf("Handle error: %v", err)
	}
	if len(summaries) != 2 {
		t.Errorf("got %d summaries; want 2", len(summaries))
	}
}
