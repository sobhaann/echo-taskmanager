package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func ConnectDB() *sql.DB {
	godotenv.Load()

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")

	int_port, _ := strconv.Atoi(port)

	psInfo := fmt.Sprintf("host=%s port=%d user=%s "+("password=%s dbname=%s sslmode=disable"),
		host, int_port, user, password, dbName)

	db, err := sql.Open("postgres", psInfo)

	if err != nil {
		log.Fatalf("there is an error in opening the database: %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("there is an error in error in connecting to the database: %v", err)
	}

	// Create tasks table if not exists
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			completed BOOLEAN DEFAULT false
		);
	`)

	if err != nil {
		log.Fatalf("there is an error in creating the tasks table: %v", err)
	}
	return db
}
