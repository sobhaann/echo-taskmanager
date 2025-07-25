package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sobhaann/echo-taskmanager/handlers"
	"github.com/sobhaann/echo-taskmanager/storage"
)

func main() {
	db := storage.ConnectDB()
	defer db.Close()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	taskHandler := &handlers.TaskHandler{DB: db}

	e.GET("/tasks", taskHandler.GetTasks)
	e.POST("/tasks", taskHandler.CreateTask)
	e.PUT("/tasks/:id", taskHandler.UpdataTask)

	e.Logger.Fatal(e.Start(":4545"))
}
