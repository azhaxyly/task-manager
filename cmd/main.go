package cmd

import (
	"fmt"
	"task-manager/internal/domain"
	"time"
)

func Run() {
	task := domain.NewTask("xaxaxax")
	if err := task.Start(); err != nil {
		fmt.Printf("Error starting task: %v\n", err)
		return
	}
	fmt.Printf("Task created: ID=%s, Status=%s, CreatedAt=%s\n", task.ID, task.Status, task.CreatedAt.Format(time.RFC3339))
}
