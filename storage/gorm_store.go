package storage

import (
	"context"

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

func (g *GormDB) CreateTask(task *models.Task, ctx context.Context) error {
	err := g.q.Task.WithContext(ctx).Create(task)
	return err
}

// i found this error that when i update a task complete status will set to false by default unless you specified to be true
func (g *GormDB) CompleteTask(id int, ctx context.Context) error {
	_, err := g.q.Task.WithContext(ctx).Where(g.q.Task.ID.Eq(id)).Update(g.q.Task.Completed, true)
	return err
}

func (g *GormDB) DeleteTask(id int, ctx context.Context) error {
	_, err := g.q.Task.WithContext(ctx).Where(g.q.Task.ID.Eq(id)).Delete()
	return err
}

func (g *GormDB) GetTasks(ctx context.Context) ([]*models.Task, error) {
	return g.q.Task.WithContext(ctx).Find()
}

func (g *GormDB) UpdateTask(id int, new_task *models.Task, ctx context.Context) error {
	current_task, err := g.q.Task.WithContext(ctx).Where(g.q.Task.ID.Eq(id)).First()
	if err != nil {
		return err
	}

	if new_task.Title == "" {
		new_task.Title = current_task.Title
	}

	if new_task.Deadline.IsZero() {
		new_task.Deadline = current_task.Deadline
	}

	new_task.CreatedAt = current_task.CreatedAt
	new_task.ID = current_task.ID

	_, err = g.q.Task.WithContext(ctx).Where(g.q.Task.ID.Eq(id)).Updates(new_task)
	if err != nil {
		return err
	}
	return nil
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
	err := g.q.User.WithContext(ctx).Create(user)
	return err
}

func (g *GormDB) GetUserByPhoneNumebr(phoneNumebr string, ctx context.Context) (*models.User, error) {
	user, err := g.q.User.WithContext(ctx).Where(g.q.User.PhoneNumber.Eq(phoneNumebr)).First()
	return user, err
}

func (g *GormDB) GetUsers(ctx context.Context) ([]*models.User, error) {
	return g.q.User.WithContext(ctx).Find()
}
