package storage

import (
	"context"

	"github.com/sobhaann/echo-taskmanager/models"
)

type Store interface {
	CreateTask(task *models.Task, ctx context.Context) error
	CompleteTask(id int, ctx context.Context) error
	DeleteTask(id int, ctx context.Context) error
	GetTasks(ctx context.Context) ([]*models.Task, error)
	UpdateTask(id int, new_task *models.Task, ctx context.Context) error
	Close() error
}

// CreateTask(task *models.Task) error
// CompleteTask(id int) error
// DeleteTask(id int) error
// GetTasks() ([]models.Task, error)
// UpdateTask(id int, new_task *models.Task) error
// Close() error
