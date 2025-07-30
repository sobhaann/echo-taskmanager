package main

import (
	"github.com/sobhaann/echo-taskmanager/handlers"
	"github.com/sobhaann/echo-taskmanager/storage"

	_ "github.com/sobhaann/echo-taskmanager/docs"
)

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
