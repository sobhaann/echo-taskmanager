package main

import (
	"github.com/sobhaann/echo-taskmanager/handlers"
	"github.com/sobhaann/echo-taskmanager/storage"

	_ "github.com/sobhaann/echo-taskmanager/docs"
)

//	@title			Echo Task Manager API
//	@version		0.01
//	@description	API for managing tasks with Echo and PostgreSQL
//	@host			localhost:4545
//	@BasePath		/

func main() {
	db := storage.ConnectDB()
	defer db.DB.Close()
	postgres := handlers.NewTaskHandler(db)
	postgres.Run()

}

//swager https://github.com/swaggo/swag
//move database stuff to storage folder
//made main.go cleaner
//openapi
