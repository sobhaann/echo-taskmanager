package handlers

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/echo/v4"
)

func (p *TaskHandler) Run() {
	//load port from `.env` file
	godotenv.Load()
	envPort := os.Getenv("PORT")
	port := fmt.Sprintf(":%s", envPort)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/tasks", p.GetTasks)
	e.POST("/tasks", p.CreateTask)
	e.PUT("/tasks/:id", p.UpdataTask)
	e.PUT("/tasks/:id/complete", p.CompleteTask)
	e.DELETE("/tasks/:id", p.DeleteTask)

	e.Logger.Fatal(e.Start(port))
}
