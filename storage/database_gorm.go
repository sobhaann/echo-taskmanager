package storage

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sobhaann/echo-taskmanager/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormStorageInterface interface {
	GormCreateTask(task *models.Task) error
	GormCompleteTask(id int) error
	GormDeleteTask(id int) error
	GormGetTasks() ([]models.Task, error)
	GormUpdateTask(id int, new_task *models.Task) error
}

type GormStore struct {
	DB *gorm.DB
}

func ConnectPostgres() *GormStore {
	godotenv.Load()

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")

	int_port, _ := strconv.Atoi(port)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		host, user, password, dbName, int_port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("there is an error in oppening the databse: %v", err)
	}

	err = db.AutoMigrate(&models.Task{})
	if err != nil {
		log.Fatalf("there is an error in creating task tabels: %v", err)
	}

	return &GormStore{
		DB: db,
	}
}

func (g *GormStore) GormCreateTask(task *models.Task) error {
	if err := g.DB.Create(task).Error; err != nil {
		return err
	}
	return nil
}

// i found this error that when i update a task complete status will set to false by default unless you specified to be true
func (g *GormStore) GormCompleteTask(id int) error {
	err := g.DB.Model(&models.Task{}).Where("id = ?", id).Update("completed", true).Error
	if err != nil {
		return err
	}
	return nil
}

func (g *GormStore) GormDeleteTask(id int) error {
	if err := g.DB.Delete(&models.Task{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (g *GormStore) GormGetTasks() ([]models.Task, error) {
	var tasks []models.Task

	result := g.DB.Find(&tasks)

	return tasks, result.Error
}

func (g *GormStore) GormUpdateTask(id int, new_task *models.Task) error {
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

	if err := g.DB.Save(new_task).Error; err != nil {
		return err
	}

	return nil
}
