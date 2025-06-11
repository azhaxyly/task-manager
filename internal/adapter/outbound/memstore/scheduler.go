package memstore

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"task-manager/internal/application/port/out"
	"task-manager/internal/common/logger"
	"task-manager/internal/domain"
)

type TaskScheduler struct {
	rng         *rand.Rand
	repo        out.TaskRepository
	mu          sync.Mutex
	cancelFuncs map[domain.TaskID]context.CancelFunc
}

func NewTaskScheduler(repo out.TaskRepository) *TaskScheduler {
	return &TaskScheduler{
		rng:         rand.New(rand.NewSource(time.Now().UnixNano())),
		repo:        repo,
		cancelFuncs: make(map[domain.TaskID]context.CancelFunc),
	}
}

func (s *TaskScheduler) Schedule(_ context.Context, id domain.TaskID) {
	ctxTask, cancel := context.WithCancel(context.Background())

	s.mu.Lock()
	s.cancelFuncs[id] = cancel
	s.mu.Unlock()

	go func() {
		task, err := s.repo.Find(ctxTask, id)
		if err != nil {
			logger.Error("task not found", "id", id, "error", err)
			return
		}
		if err := task.Start(); err != nil {
			logger.Error("failed to start task", "id", id, "error", err)
			return
		}
		s.repo.Save(ctxTask, task)

		// имитация задержки
		delay := time.Duration(s.rng.Intn(3))*time.Minute + 3*time.Minute

		select {
		case <-time.After(delay):
			t2, err := s.repo.Find(ctxTask, id)
			if err != nil {
				logger.Error("task not found after delay", "id", id, "error", err)
				return
			}
			_ = t2.Complete(string(domain.Success))
			s.repo.Save(ctxTask, t2)

		case <-ctxTask.Done():
			t2, err := s.repo.Find(context.Background(), id)
			if err != nil {
				logger.Error("task not found on cancel", "id", id, "error", err)
				return
			}
			_ = t2.Cancel()
			s.repo.Save(context.Background(), t2)
		}

		s.mu.Lock()
		delete(s.cancelFuncs, id)
		s.mu.Unlock()
	}()
}

func (s *TaskScheduler) Cancel(_ context.Context, id domain.TaskID) {
	s.mu.Lock()
	cancel, ok := s.cancelFuncs[id]
	s.mu.Unlock()

	if ok {
		cancel()
	}
}
