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
func (thg *TaskHandlerGORM) GormRun() {
	//load port from `.env` file
	godotenv.Load()
	envPort := os.Getenv("PORT")
	port := fmt.Sprintf(":%s", envPort)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//swagger ui
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/tasks", thg.GetTasks)
	e.POST("/tasks", thg.CreateTask)
	e.PUT("/tasks/:id", thg.UpdataTask)
	e.PUT("/tasks/:id/complete", thg.CompleteTask)
	e.DELETE("/tasks/:id", thg.DeleteTask)

	e.Logger.Fatal(e.Start(port))
}

func (thpq *TaskHandlerPQ) PQRun() {
	//load port from `.env` file
	godotenv.Load()
	envPort := os.Getenv("PORT")
	port := fmt.Sprintf(":%s", envPort)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//swagger ui
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/tasks", thpq.GetTasks)
	e.POST("/tasks", thpq.CreateTask)
	e.PUT("/tasks/:id", thpq.UpdataTask)
	e.PUT("/tasks/:id/complete", thpq.CompleteTask)
	e.DELETE("/tasks/:id", thpq.DeleteTask)

	e.Logger.Fatal(e.Start(port))
}
