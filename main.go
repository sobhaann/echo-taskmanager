package main

import (
	"log"

	_ "github.com/sobhaann/echo-taskmanager/docs"
	"github.com/sobhaann/echo-taskmanager/handlers"
	"github.com/sobhaann/echo-taskmanager/storage"
)

func main() {
	db, err := storage.InitDB()
	if err != nil {
		log.Fatalf("there is an error in initialize the db: %v", err)
	}
	handler := handlers.NewHandler(db)
	handler.Run()

}

// --TODO--
//swager https://github.com/swaggo/swag
//move database stuff to storage folder
//made main.go cleaner
//openapi
