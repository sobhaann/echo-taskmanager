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

	query := `INSERT INTO tasks (title) VALUES ($1) RETURNING id`

	err := h.DB.QueryRow(query, task.Title).Scan(&task)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, task)
}

// get tasks from db and return it
func (h *TaskHandler) GetTasks(c echo.Context) error {
	query := `SELECT id, title, completed FROM tasks`
	rows, err := h.DB.Query(query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()

	tasks := []models.Task{}
	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Completed)
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
	id, _ := strconv.Atoi(c.Param("id"))
	task := new(models.Task)
	if err := c.Bind(task); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	update_query := `UPDATE tasks SET title = $1 WHERE id = $2`
	_, err := h.DB.Exec(update_query, task.Title, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	task.ID = id
	return c.JSON(http.StatusOK, map[string]models.Task{"updated task": *task})
}
