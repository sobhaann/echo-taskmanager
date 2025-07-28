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
	//load port from `.env` file
	godotenv.Load()
	envPort := os.Getenv("PORT")
	port := fmt.Sprintf(":%s", envPort)

	db := storage.ConnectDB()
	defer db.DB.Close()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	postgres := handlers.NewTaskHandler(db)

	e.GET("/tasks", postgres.GetTasks)
	e.POST("/tasks", postgres.CreateTask)
	e.PUT("/tasks/:id", postgres.UpdataTask)
	e.PUT("/tasks/:id/complete", postgres.CompleteTask)
	e.DELETE("/tasks/:id", postgres.DeleteTask)

	e.Logger.Fatal(e.Start(port))
}

//swager https://github.com/swaggo/swag
//move database stuff to storage folder
//made main.go cleaner
//openapi
