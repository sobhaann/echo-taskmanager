package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sobhaann/echo-taskmanager/models"
	"github.com/sobhaann/echo-taskmanager/storage"
)

type TaskHandler struct {
	store storage.StorageInterface
}

func NewTaskHandler(store storage.StorageInterface) *TaskHandler {
	return &TaskHandler{
		store: store,
	}
}

// create a new task
func (h *TaskHandler) CreateTask(c echo.Context) error {
	task := new(models.Task)
	if err := c.Bind(task); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	//checking if any of the values is zero or not; if any value is a zero value it return an error

	if task.Deadline.IsZero() {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Deadline is required"})
	}

	if task.Title == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "title is requierd"})
	}

	err := h.store.DBCreateTask(task)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, task)
}

// get tasks from db and return it
func (h *TaskHandler) GetTasks(c echo.Context) error {
	// query := `SELECT id, title, completed, created_at, deadline FROM tasks`
	// rows, err := h.DB.Query(query)
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, err.Error())
	// }
	// defer rows.Close()

	// tasks := []models.Task{}
	// for rows.Next() {
	// 	var task models.Task
	// 	err := rows.Scan(&task.ID, &task.Title, &task.Completed, &task.CreatedAt, &task.Deadline)
	// 	if err != nil {
	// 		return c.JSON(http.StatusInternalServerError, err.Error())
	// 	}

	// 	tasks = append(tasks, task)
	// }
	// if err = rows.Err(); err != nil {
	// 	return c.JSON(http.StatusInternalServerError, err.Error())
	// }

	// return c.JSON(http.StatusOK, tasks)
	tasks, err := h.store.DBGetTasks()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, tasks)
}

// update tasks
func (h *TaskHandler) UpdataTask(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	new_task := new(models.Task)
	if err := c.Bind(new_task); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if new_task.Deadline.IsZero() && new_task.Title == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "all of the fields are empty"})
	}

	err := h.store.DBUpdateTask(id, new_task)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	new_task.ID = id
	return c.JSON(http.StatusOK, map[string]models.Task{"updated task": *new_task})
}

// complete a task
func (h *TaskHandler) CompleteTask(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	err := h.store.DBCompleteTask(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Task completed"})
}

// delete a task
func (h *TaskHandler) DeleteTask(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	err := h.store.DBDeleteTask(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "task deleted"})
}
