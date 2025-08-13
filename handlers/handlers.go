package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	_ "github.com/sobhaann/echo-taskmanager/docs"
	"github.com/sobhaann/echo-taskmanager/models"

	"github.com/sobhaann/echo-taskmanager/storage"
)

type Handler struct {
	store storage.Store
}

func NewHandler(store storage.Store) *Handler {
	return &Handler{
		store: store,
	}
}

// CreateTask godoc
//
//	@Summary		Create task
//	@Description	create a new task
//	@Tags			tasks
//	@Accept			json
//	@Produce		json
//	@Param			task	body		models.Task	true	"Task object"
//	@Success		201		{object}	models.Task
//	@Router			/tasks [post]
func (h *Handler) CreateTask(c echo.Context) error {
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
	user_id := h.claimUserID(c)
	err := h.store.CreateTask(c.Request().Context(), task, user_id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, task)
}

// GetTasks godoc
//
//	@Summary		List tasks
//	@Description	get all tasks from database
//	@Tags			tasks
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	models.Task
//	@Router			/tasks [get]
func (h *Handler) GetTasks(c echo.Context) error {
	user_id := h.claimUserID(c)
	tasks, err := h.store.GetTasks(c.Request().Context(), user_id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, tasks)
}

// UpdataTask godoc
//
//	@Summary		Update task
//	@Description	update an existing task
//	@Tags			tasks
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Task ID"
//
//	@Success		200	{object}	models.Task
//	@Router			/tasks/{id} [put]
func (h *Handler) UpdataTask(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	new_task := new(models.Task)
	if err := c.Bind(new_task); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if new_task.Deadline.IsZero() && new_task.Title == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "all of the fields are empty"})
	}
	user_id := h.claimUserID(c)
	err := h.store.UpdateTask(c.Request().Context(), new_task, id, user_id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	new_task.ID = id
	return c.JSON(http.StatusOK, map[string]models.Task{"updated task": *new_task})
}

// CompleteTask godoc
//
//	@Summary		Complete task
//	@Description	mark a task as completed
//	@Tags			tasks
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Task ID"
//	@Success		200	{object}	models.Task
//	@Router			/tasks/{id}/complete [put]
func (h *Handler) CompleteTask(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	user_id := h.claimUserID(c)
	err := h.store.CompleteTask(c.Request().Context(), id, user_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Task completed"})
}

// DeleteTask godoc
//
//	@Summary		Delete task
//	@Description	delete a task
//	@Tags			tasks
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Task ID"
//	@Success		200	{object}	nil
//	@Router			/tasks/{id} [delete]
func (h *Handler) DeleteTask(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	user_id := h.claimUserID(c)
	err := h.store.DeleteTask(c.Request().Context(), id, user_id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "task deleted"})
}
