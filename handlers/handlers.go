package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	_ "github.com/sobhaann/echo-taskmanager/docs"
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

// CreateTask godoc
// @Summary      Create a new task
// @Description  Create a new task with the provided details
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        task  body  models.Task  true  "Task to create"
// @Success      201   {object}  models.Task
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /tasks [post]

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

// GetTasks godoc
// @Summary      Get all tasks
// @Description  Retrieve all tasks from the database
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Success      200  {array}  models.Task
// @Failure      500  {object}  map[string]string
// @Router       /tasks [get]

// get tasks from db and return it
func (h *TaskHandler) GetTasks(c echo.Context) error {
	tasks, err := h.store.DBGetTasks()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, tasks)
}

// UpdataTask godoc
// @Summary      Update a task
// @Description  Update a task's details by its ID
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id    path   int          true  "Task ID"
// @Param        task  body   models.Task  true  "Task updates"
// @Success      200   {object}  models.Task
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /tasks/{id} [put]

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

// CompleteTask godoc
// @Summary      Mark a task as completed
// @Description  Mark a task as completed by its ID
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Task ID"
// @Success      200  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /tasks/{id}/complete [put]

// complete a task
func (h *TaskHandler) CompleteTask(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	err := h.store.DBCompleteTask(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Task completed"})
}

// DeleteTask godoc
// @Summary      Delete a task
// @Description  Delete a task by its ID
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Task ID"
// @Success      200  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /tasks/{id} [delete]

// delete a task
func (h *TaskHandler) DeleteTask(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	err := h.store.DBDeleteTask(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "task deleted"})
}
