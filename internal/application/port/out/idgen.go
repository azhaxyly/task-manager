package out

import "task-manager/internal/domain"

type IDGenerator interface {
	NewID() domain.TaskID
}
