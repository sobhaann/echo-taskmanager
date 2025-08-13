package handlers

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
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
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func (h *Handler) Run() {
	// Load environment variables
	godotenv.Load()
	envPort := os.Getenv("PORT")
	port := fmt.Sprintf(":%s", envPort)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// ===== Public routes =====
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.POST("/signup", h.Signup)
	e.POST("/login", h.Login)
	e.GET("/users", h.GetUsers) // If you want to make it public

	e.GET("/tasks", h.GetTasks)
	e.POST("/tasks", h.CreateTask)
	e.PUT("/tasks/:id", h.UpdataTask)
	e.PUT("/tasks/:id/complete", h.CompleteTask)
	e.DELETE("/tasks/:id", h.DeleteTask)

	// ===== Private routes =====
	secret := os.Getenv("JWT_SECRET")
	r := e.Group("")
	r.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(secret),
		//TokenLookup: "header:Authorization", // it doesnt work very well
		ContextKey: "user",
	}))
	r.GET("/profile", h.Profile)

	defer h.store.Close()
	e.Logger.Fatal(e.Start(port))
}
