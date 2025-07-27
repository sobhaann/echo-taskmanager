package main

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sobhaann/echo-taskmanager/handlers"
	"github.com/sobhaann/echo-taskmanager/storage"

	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()
	envPort := os.Getenv("PORT")
	port := fmt.Sprintf(":%s", envPort)
	db := storage.ConnectDB()
	defer db.Close()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	taskHandler := &handlers.TaskHandler{DB: db}

	e.GET("/tasks", taskHandler.GetTasks)
	e.POST("/tasks", taskHandler.CreateTask)
	e.PUT("/tasks/:id", taskHandler.UpdataTask)
	e.PUT("/tasks/:id/complete", taskHandler.CompleteTask)
	e.DELETE("/tasks/:id", taskHandler.DeleteTask)

	e.Logger.Fatal(e.Start(port))
}
