package domain

import "fmt"

type Status string

const (
	// таск создан, но пока не запущен
	Pending Status = "pending"
	// таск запущен и выполняется
	Running Status = "running"
	// таск завершился успешно
	Success Status = "completed"
	// таск завершился с ошибкой
	Failed Status = "failed"
	// таск отменен
	Canceled Status = "canceled"
)

var AllStatuses = []Status{
	Pending,
	Running,
	Success,
	Failed,
	Canceled,
}

func (s Status) IsValid() bool {
	switch s {
	case Pending, Running, Success, Failed, Canceled:
		return true
	}
	return false
}

func ParseStatus(raw string) (Status, error) {
	s := Status(raw)
	if !s.IsValid() {
		return "", fmt.Errorf("invalid status value: %q", raw)
	}
	return s, nil
}

func (s Status) IsTerminal() bool {
	return s == Success || s == Failed || s == Canceled
}
