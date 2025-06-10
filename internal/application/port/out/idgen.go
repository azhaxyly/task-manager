package out

import "task-manager/internal/domain"

type IDGenerator interface {
	NewID() domain.TaskID
}

func NewStaticID(id string) IDGenerator {
	return staticID(id)
}

type staticID string

func (s staticID) NewID() domain.TaskID {
	return domain.TaskID(s)
}
