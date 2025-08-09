package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sobhaann/echo-taskmanager/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormDB struct {
	DB *gorm.DB
}

type PqDB struct {
	DB *sql.DB
}

func LoadDatabaseInfo() (string, string) {
	godotenv.Load()

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbEngine := os.Getenv("DB_ENGINE")

	int_port, _ := strconv.Atoi(port)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		host, user, password, dbName, int_port)

	return dsn, strings.ToLower(dbEngine)
}

func InitDB() (Store, error) {
	dsn, dbEngine := LoadDatabaseInfo()

	switch dbEngine {
	case "gorm":
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("there is an error in oppening the databse: %v", err)
		}

		err = db.AutoMigrate(&models.Task{})
		if err != nil {
			log.Fatalf("there is an error in creating task tabels: %v", err)
		}
		return &GormDB{
			DB: db,
		}, nil
	case "pq":
		db, err := sql.Open("postgres", dsn)
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
			completed BOOLEAN DEFAULT false,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			deadline TIMESTAMP WITH TIME ZONE
		);
	`)
		if err != nil {
			log.Fatalf("there is an error in creating the tasks table: %v", err)
		}
		return &PqDB{
			DB: db,
		}, nil

	default:
		return nil, fmt.Errorf("unsupported DB engine: %s", dbEngine)
	}
}
