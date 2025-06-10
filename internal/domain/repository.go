package domain

import "context"

type TaskRepository interface {
	// сохраняет либо обновляет таску
	Save(ctx context.Context, task *Task) error

	// ищет таску по айдишке
	Find(ctx context.Context, id TaskID) (*Task, error)

	// удаялет таску
	Delete(ctx context.Context, id TaskID) error
}
