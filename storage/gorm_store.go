package storage

import (
	"github.com/sobhaann/echo-taskmanager/models"
)

func (g *GormDB) CreateTask(task *models.Task) error {
	if err := g.DB.Create(task).Error; err != nil {
		return err
	}
	return nil
}

// i found this error that when i update a task complete status will set to false by default unless you specified to be true
func (g *GormDB) CompleteTask(id int) error {
	err := g.DB.Model(&models.Task{}).Where("id = ?", id).Update("completed", true).Error
	if err != nil {
		return err
	}
	return nil
}

func (g *GormDB) DeleteTask(id int) error {
	if err := g.DB.Delete(&models.Task{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (g *GormDB) GetTasks() ([]models.Task, error) {
	var tasks []models.Task

	result := g.DB.Find(&tasks)

	return tasks, result.Error
}

func (g *GormDB) UpdateTask(id int, new_task *models.Task) error {
	var current_task models.Task
	if err := g.DB.First(&current_task, id).Error; err != nil {
		return err
	}

	if current_task.Title == "" {
		new_task.Title = current_task.Title
	}

	if new_task.Deadline.IsZero() {
		new_task.Deadline = current_task.Deadline
	}

	new_task.ID = current_task.ID
	new_task.CreatedAt = current_task.CreatedAt

	if err := g.DB.Save(new_task).Error; err != nil {
		return err
	}

	return nil
}

func (g *GormDB) Close() error {
	db, err := g.DB.DB()
	if err != nil {
		return err
	}
	return db.Close()
}
