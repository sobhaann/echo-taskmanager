package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sobhaann/echo-taskmanager/models"
)

type PostgresStorageInterface interface {
	PQCreateTask(t *models.Task) error
	PQCompleteTask(id int) error
	PQDeleteTask(id int) error
	PQGetTasks() ([]models.Task, error)
	PQUpdateTask(id int, new_task *models.Task) error
}

type PostgresStore struct {
	DB *sql.DB
}

// creating the database schema and connect to it
func ConnectPostgresPQ() *PostgresStore {
	godotenv.Load()

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")

	int_port, _ := strconv.Atoi(port)

	psInfo := fmt.Sprintf("host=%s port=%d user=%s "+("password=%s dbname=%s sslmode=disable"),
		host, int_port, user, password, dbName)

	db, err := sql.Open("postgres", psInfo)

	if err != nil {
		log.Fatalf("there is an error in opening the database: %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("there is an error in error in connecting to the database: %v", err)
	}

	// Create tasks table if not exists
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			completed BOOLEAN DEFAULT false,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			deadline TIMESTAMP WITH TIME ZONE
		);
	`)

	if err != nil {
		log.Fatalf("there is an error in creating the tasks table: %v", err)
	}
	return &PostgresStore{
		DB: db,
	}
}

func (p *PostgresStore) PQCreateTask(t *models.Task) error {
	query := `INSERT INTO tasks (title,created_at, deadline) VALUES ($1, CURRENT_TIMESTAMP, $2) RETURNING id, created_at`
	err := p.DB.QueryRow(query, t.Title, t.Deadline).Scan(&t.ID, &t.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresStore) PQCompleteTask(id int) error {
	query := `UPDATE tasks SET completed = true WHERE id = $1`
	_, err := p.DB.Exec(query, id)
	return err
}

func (p *PostgresStore) PQDeleteTask(id int) error {
	query := `DELETE FROM tasks WHERE ID = $1`
	_, err := p.DB.Exec(query, id)
	return err
}

func (p *PostgresStore) PQGetTasks() ([]models.Task, error) {
	query := `SELECT * FROM tasks`
	rows, err := p.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []models.Task{}
	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Completed, &task.CreatedAt, &task.Deadline)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (p *PostgresStore) PQGetTaskByID(id int) (*models.Task, error) {
	current_task := &models.Task{}
	query := `SELECT title, created_at, deadline FROM tasks WHERE id = $1`
	err := p.DB.QueryRow(query, id).Scan(&current_task.Title, &current_task.CreatedAt, &current_task.Deadline)
	if err != nil {
		return current_task, err
	}
	return current_task, nil
}

func (p *PostgresStore) PQUpdateTask(id int, new_task *models.Task) error {
	current_task, err := p.PQGetTaskByID(id)
	if err != nil {
		return err
	}

	if new_task.Title == "" {
		new_task.Title = current_task.Title
	}
	if new_task.Deadline.IsZero() {
		new_task.Deadline = current_task.Deadline
	}

	update_query := `UPDATE tasks SET title = $1, deadline = $2 WHERE id = $3`
	_, err = p.DB.Exec(update_query, new_task.Title, new_task.Deadline, id)
	if err != nil {
		return err
	}
	return nil
}
