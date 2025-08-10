package storage

import (
	_ "github.com/lib/pq"
	"github.com/sobhaann/echo-taskmanager/models"
)

func (p *PqDB) CreateTask(t *models.Task) error {
	query := `INSERT INTO tasks (title,created_at, deadline) VALUES ($1, CURRENT_TIMESTAMP, $2) RETURNING id, created_at`
	err := p.DB.QueryRow(query, t.Title, t.Deadline).Scan(&t.ID, &t.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (p *PqDB) CompleteTask(id int) error {
	query := `UPDATE tasks SET completed = true WHERE id = $1`
	_, err := p.DB.Exec(query, id)
	return err
}

func (p *PqDB) DeleteTask(id int) error {
	query := `DELETE FROM tasks WHERE ID = $1`
	_, err := p.DB.Exec(query, id)
	return err
}

func (p *PqDB) GetTasks() ([]models.Task, error) {
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

func (p *PqDB) PQGetTaskByID(id int) (*models.Task, error) {
	current_task := &models.Task{}
	query := `SELECT title, created_at, deadline FROM tasks WHERE id = $1`
	err := p.DB.QueryRow(query, id).Scan(&current_task.Title, &current_task.CreatedAt, &current_task.Deadline)
	if err != nil {
		return current_task, err
	}
	return current_task, nil
}

func (p *PqDB) UpdateTask(id int, new_task *models.Task) error {
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

func (p *PqDB) Close() error {
	return p.DB.Close()
}
