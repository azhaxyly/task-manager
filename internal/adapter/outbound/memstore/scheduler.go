package memstore

import (
	"context"
	"sync"
	"task-manager/internal/application/port/out"
	"task-manager/internal/domain"
)

type TaskScheduler struct {
	repo        out.TaskRepository
	mu          sync.Mutex
	cancelFuncs map[domain.TaskID]context.CancelFunc
}

func NewTaskScheduler(repo out.TaskRepository) *TaskScheduler {
	return &TaskScheduler{
		repo:        repo,
		cancelFuncs: make(map[domain.TaskID]context.CancelFunc),
	}
}

func (s *TaskScheduler) Schedule(_ context.Context, id domain.TaskID) {
	_, cancel := context.WithCancel(context.Background())

	s.mu.Lock()
	s.cancelFuncs[id] = cancel
	s.mu.Unlock()

	// TODO: заебашить горутину
}
