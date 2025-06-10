package memstore_test

import (
	"context"
	"testing"

	"task-manager/internal/adapter/outbound/memstore"
	"task-manager/internal/domain"
)

func TestTaskRepository_SaveFindDeleteList(t *testing.T) {
	ctx := context.Background()
	repo := memstore.NewTaskRepository()

	id := domain.TaskID("test-id")
	orig := domain.NewTask(id)
	if err := repo.Save(ctx, orig); err != nil {
		t.Fatalf("Save error: %v", err)
	}

	got, err := repo.Find(ctx, id)
	if err != nil {
		t.Fatalf("Find error: %v", err)
	}
	if got.ID != id || got.Status != domain.Pending {
		t.Errorf("got %+v, want ID %s status %s", got, id, domain.Pending)
	}

	got.Status = domain.Running
	again, _ := repo.Find(ctx, id)
	if again.Status != domain.Pending {
		t.Errorf("stored task was mutated; want status %s, got %s", domain.Pending, again.Status)
	}

	ids, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if len(ids) != 1 || ids[0] != id {
		t.Errorf("List returned %v; want [%s]", ids, id)
	}

	if err := repo.Delete(ctx, id); err != nil {
		t.Fatalf("Delete error: %v", err)
	}
	if _, err := repo.Find(ctx, id); err != domain.ErrTaskNotFound {
		t.Errorf("after delete, Find should ErrTaskNotFound; got %v", err)
	}

	if err := repo.Delete(ctx, "nope"); err != domain.ErrTaskNotFound {
		t.Errorf("Delete non-existent should ErrTaskNotFound; got %v", err)
	}
}
