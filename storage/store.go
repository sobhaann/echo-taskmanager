package storage

import (
	"context"

	"github.com/sobhaann/echo-taskmanager/models"
)

type Store interface {
	TaskStore
	UserStore
}

type TaskStore interface {
	CreateTask(task *models.Task, ctx context.Context) error
	CompleteTask(id int, ctx context.Context) error
	DeleteTask(id int, ctx context.Context) error
	GetTasks(ctx context.Context) ([]*models.Task, error)
	UpdateTask(id int, new_task *models.Task, ctx context.Context) error
	Close() error
}

type UserStore interface {
	CreateUser(user *models.User, ctx context.Context) error
	GetUserByPhoneNumebr(phone_number string, ctx context.Context) (*models.User, error)
	GetUsers(ctx context.Context) ([]*models.User, error)
}
