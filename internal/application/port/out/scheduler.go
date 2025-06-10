package out

import (
	"context"
	"task-manager/internal/domain"
)

type TaskScheduler interface {
	// запуск планировщик обработки задач с известной айдишкой
	Schedule(ctx context.Context, id domain.TaskID)
	// стопит планировщика
	Cancel(ctx context.Context, id domain.TaskID)
	// перезапуск с таймером
	Reschedule(ctx context.Context, id domain.TaskID, newTime int64)
	// гет всех запланированных тасок
	GetScheduledTasks(ctx context.Context) ([]domain.Task, error)
}
