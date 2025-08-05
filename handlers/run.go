package handlers

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/sobhaann/echo-taskmanager/docs" // Swagger generated docs
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title			Echo Task Manager API
// @version		0.01
// @description	API for managing tasks with Echo and PostgreSQL
// @host			localhost:4545
// @BasePath		/
func (th *TaskHandler) Run() {
	//load port from `.env` file
	godotenv.Load()
	envPort := os.Getenv("PORT")
	port := fmt.Sprintf(":%s", envPort)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//swagger ui
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/tasks", th.GetTasks)
	e.POST("/tasks", th.CreateTask)
	e.PUT("/tasks/:id", th.UpdataTask)
	e.PUT("/tasks/:id/complete", th.CompleteTask)
	e.DELETE("/tasks/:id", th.DeleteTask)

	e.Logger.Fatal(e.Start(port))
}
