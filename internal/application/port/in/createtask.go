package in

import "context"

// TODO: пока хз
type CreateTaskCommand struct{}

type CreateTaskUseCase interface {
	Handle(ctx context.Context, cmd CreateTaskCommand) error
	Validate(cmd CreateTaskCommand) error
}
