package storage

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/sobhaann/echo-taskmanager/models"
)

type PqDB struct {
	db *sql.DB
}

func NewPqDB(dsn string) (*PqDB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	//user Table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		user_name TEXT,
		password TEXT NOT NULL,
		phone_number TEXT NOT NULL UNIQUE
		);
		`)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			completed BOOLEAN DEFAULT false,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			deadline TIMESTAMP WITH TIME ZONE,
			user_id INT,
			CONSTRAINT fk_user FOREIGN KEY (user_id)
			REFERENCES users(id)
		);
	`)
	if err != nil {
		return nil, err
	}

	return &PqDB{db: db}, nil
}

func (p *PqDB) PQGetTaskByID(id int, user_id int, ctx context.Context) (*models.Task, error) {
	current_task := &models.Task{}
	query := `SELECT id, title, created_at, deadline, completed, user_id FROM tasks WHERE id = $1 AND user_id = $2`
	err := p.db.QueryRowContext(ctx, query, id, user_id).
		Scan(&current_task.ID, &current_task.Title, &current_task.CreatedAt, &current_task.Deadline, &current_task.Completed, &current_task.UserID)
	if err != nil {
		return current_task, err
	}
	return current_task, nil
}

func (p *PqDB) CreateTask(ctx context.Context, task *models.Task, user_id int) error {
	query := `INSERT INTO tasks (title,created_at, deadline, user_id) VALUES ($1, CURRENT_TIMESTAMP, $2, $3) RETURNING id, created_at`
	err := p.db.QueryRowContext(ctx, query, task.Title, task.Deadline, user_id).Scan(&task.ID, &task.CreatedAt)
	return err
}

func (p *PqDB) CompleteTask(ctx context.Context, id int, user_id int) error {
	query := `UPDATE tasks SET completed = true WHERE id = $1 AND user_id = $2`
	_, err := p.db.ExecContext(ctx, query, id, user_id)

	return err
}

func (p *PqDB) DeleteTask(ctx context.Context, id int, user_id int) error {
	query := `DELETE FROM tasks WHERE id = $1 AND user_id = $2`
	_, err := p.db.ExecContext(ctx, query, id, user_id)

	return err
}

func (p *PqDB) GetTasks(ctx context.Context, user_id int) ([]*models.Task, error) {
	query := `SELECT * FROM tasks WHERE user_id = $1`
	rows, err := p.db.QueryContext(ctx, query, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Completed, &task.CreatedAt, &task.Deadline, &task.UserID)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)

	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (p *PqDB) UpdateTask(ctx context.Context, new_task *models.Task, id int, user_id int) error {
	current_task, err := p.PQGetTaskByID(id, user_id, ctx)
	// if current_task.UserID == user_id {
	// 	return errors.New("invalid jwt token for update")
	// }
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
	new_task.Completed = current_task.Completed
	new_task.UserID = user_id
	new_task.ID = current_task.ID

	update_query := `UPDATE tasks SET title = $1, deadline = $2 WHERE id = $3 AND user_id = $4`
	_, err = p.db.Exec(update_query, new_task.Title, new_task.Deadline, id, user_id)
	if err != nil {
		return err
	}
	return nil
}

// auth
func (p *PqDB) CreateUser(user *models.User, ctx context.Context) error {
	query := `INSERT INTO users (user_name, password, phone_number) VALUES ($1, $2, $3) RETURNING id`
	err := p.db.QueryRowContext(ctx, query, user.UserName, user.Password, user.PhoneNumber).Scan(&user.ID)
	return err

}

func (p *PqDB) GetUserByPhoneNumber(phone_number string, ctx context.Context) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, user_name, phone_number, password FROM users WHERE phone_number = $1`
	err := p.db.QueryRow(query, phone_number).
		Scan(&user.ID, &user.UserName, &user.PhoneNumber, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (p *PqDB) GetUsers(ctx context.Context) ([]*models.User, error) {
	query := `SELECT * FROM users`
	rows, err := p.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.UserName, &user.Password, &user.PhoneNumber)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (p *PqDB) Close() error {
	return p.db.Close()
}
