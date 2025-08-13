package storage

import (
	"context"
	"errors"

	"github.com/sobhaann/echo-taskmanager/dao"
	"github.com/sobhaann/echo-taskmanager/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormDB struct {
	q  *dao.Query
	db *gorm.DB
}

func NewGormDB(dsn string) (*GormDB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&models.Task{}); err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&models.User{}); err != nil {
		return nil, err
	}
	q := dao.Use(db)
	return &GormDB{
		db: db,
		q:  q,
	}, nil
}

func (g *GormDB) CreateTask(ctx context.Context, task *models.Task, user_id int) error {
	task.UserID = user_id
	err := g.q.Task.WithContext(ctx).Create(task)
	return err
}

// i found this error that when i update a task complete status will set to false by default unless you specified to be true
func (g *GormDB) CompleteTask(ctx context.Context, id int, user_id int) error {
	task, err := g.q.Task.WithContext(ctx).Where(g.q.Task.ID.Eq(id)).First()
	if err != nil {
		return err
	}
	if task.UserID == user_id {
		_, err := g.q.Task.WithContext(ctx).Where(g.q.Task.ID.Eq(id)).Update(g.q.Task.Completed, true)
		return err
	}
	return errors.New("you are not allowed to complete this task")
}

func (g *GormDB) DeleteTask(ctx context.Context, id int, user_id int) error {
	task, err := g.q.Task.WithContext(ctx).Where(g.q.Task.ID.Eq(id)).First()
	if err != nil {
		return err
	}
	if task.UserID == user_id {
		_, err := g.q.Task.WithContext(ctx).Where(g.q.Task.ID.Eq(id)).Delete()
		return err
	}
	return errors.New("you are not allowed to delete this task")

}

func (g *GormDB) GetTasks(ctx context.Context, user_id int) ([]*models.Task, error) {
	return g.q.Task.WithContext(ctx).Where(g.q.Task.UserID.Eq(user_id)).Find()
}

func (g *GormDB) UpdateTask(ctx context.Context, new_task *models.Task, id int, user_id int) error {
	current_task, err := g.q.Task.WithContext(ctx).Where(g.q.Task.ID.Eq(id)).First()
	if err != nil {
		return err
	}

	if current_task.UserID != user_id {
		return errors.New("you are not allowed to update this task")
	}

	if new_task.Title == "" {
		new_task.Title = current_task.Title
	}

	if new_task.Deadline.IsZero() {
		new_task.Deadline = current_task.Deadline
	}

	new_task.CreatedAt = current_task.CreatedAt
	new_task.UserID = current_task.UserID
	new_task.ID = current_task.ID

	_, err = g.q.Task.WithContext(ctx).Where(g.q.Task.ID.Eq(id)).Updates(new_task)
	return err
}

func (g *GormDB) Close() error {
	db, err := g.db.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

// Authitication
func (g *GormDB) CreateUser(user *models.User, ctx context.Context) error {
	existing, _ := g.q.User.WithContext(ctx).Where(g.q.User.PhoneNumber.Eq(user.PhoneNumber)).First()
	if existing != nil {
		return errors.New("user with this phone number already exists")
	}
	err := g.q.User.WithContext(ctx).Create(user)
	return err
}

func (g *GormDB) GetUserByPhoneNumber(phoneNumber string, ctx context.Context) (*models.User, error) {
	user, err := g.q.User.WithContext(ctx).Where(g.q.User.PhoneNumber.Eq(phoneNumber)).First()
	return user, err
}

func (g *GormDB) GetUsers(ctx context.Context) ([]*models.User, error) {
	return g.q.User.WithContext(ctx).Find()
}
