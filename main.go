package main

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sobhaann/echo-taskmanager/handlers"
	"github.com/sobhaann/echo-taskmanager/storage"

	_ "github.com/sobhaann/echo-taskmanager/docs"
)

func main() {
	godotenv.Load()
	dbEngineType := os.Getenv("DB_ENGINE")

	if strings.ToLower(dbEngineType) == "gorm" {
		db := storage.ConnectPostgresGORM()

		postgres := handlers.NewTaskHandlerGORM(db)

		postgres.GormRun()
	} else if strings.ToLower(dbEngineType) == "pq" {
		db := storage.ConnectPostgresPQ()

		postgres := handlers.NewTaskHandlerPQ(db)

		postgres.PQRun()
	} else {
		log.Fatal("the database engine should be `gorm` or `pq`!!!!")
	}

}

// --TODO--
//swager https://github.com/swaggo/swag
//move database stuff to storage folder
//made main.go cleaner
//openapi
