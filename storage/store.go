package storage

import "github.com/sobhaann/echo-taskmanager/models"

type Store interface {
	CreateTask(task *models.Task) error
	CompleteTask(id int) error
	DeleteTask(id int) error
	GetTasks() ([]models.Task, error)
	UpdateTask(id int, new_task *models.Task) error
	Close() error
}
