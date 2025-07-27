package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sobhaann/echo-taskmanager/models"
)

type TaskHandler struct {
	DB *sql.DB
}

// create a new task
func (h *TaskHandler) CreateTask(c echo.Context) error {
	task := new(models.Task)
	if err := c.Bind(task); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if task.Deadline.IsZero() {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Deadline is required"})
	}

	query := `INSERT INTO tasks (title,created_at, deadline) VALUES ($1, CURRENT_TIMESTAMP, $2) RETURNING id, created_at`

	err := h.DB.QueryRow(query, task.Title, task.Deadline).Scan(&task.ID, &task.CreatedAt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, task)
}

// get tasks from db and return it
func (h *TaskHandler) GetTasks(c echo.Context) error {
	query := `SELECT id, title, completed, created_at, deadline FROM tasks`
	rows, err := h.DB.Query(query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()

	tasks := []models.Task{}
	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Completed, &task.CreatedAt, &task.Deadline)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, tasks)
}

// update tasks
func (h *TaskHandler) UpdataTask(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	task := new(models.Task)
	if err := c.Bind(task); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if task.Deadline.IsZero() {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Deadline is required"})
	}

	update_query := `UPDATE tasks SET title = $1, deadline = $2 WHERE id = $3`
	_, err := h.DB.Exec(update_query, task.Title, task.Deadline, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	task.ID = id
	return c.JSON(http.StatusOK, map[string]models.Task{"updated task": *task})
}

// complete a task
func (h *TaskHandler) CompleteTask(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	query := `UPDATE tasks SET completed = true WHERE id = $1`
	_, err := h.DB.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Task completed"})
}

// delete a task
func (h *TaskHandler) DeleteTask(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	query := `DELETE FROM tasks WHERE ID = $1`
	_, err := h.DB.Exec(query, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "task deleted"})
}
