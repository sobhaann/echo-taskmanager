package storage

import (
	"context"

	"github.com/sobhaann/echo-taskmanager/models"
)

type Store interface {
	TaskStore
	UserStore
	Close() error
}

type TaskStore interface {
	CreateTask(ctx context.Context, task *models.Task, user_id int) error
	CompleteTask(ctx context.Context, id int, user_id int) error
	DeleteTask(ctx context.Context, id int, user_id int) error
	GetTasks(ctx context.Context, user_id int) ([]*models.Task, error)
	UpdateTask(ctx context.Context, new_task *models.Task, id int, user_id int) error
}

type UserStore interface {
	CreateUser(user *models.User, ctx context.Context) error
	GetUserByPhoneNumber(phone_number string, ctx context.Context) (*models.User, error)
	GetUsers(ctx context.Context) ([]*models.User, error)
}
