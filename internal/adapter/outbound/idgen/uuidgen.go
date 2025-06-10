package idgen

import (
	"task-manager/internal/domain"

	"github.com/google/uuid"
)

type UUIDGenerator struct{}

func NewUUIDGenerator() *UUIDGenerator {
	return &UUIDGenerator{}
}

func (g *UUIDGenerator) NewID() domain.TaskID {
	return domain.TaskID(uuid.NewString())
}
