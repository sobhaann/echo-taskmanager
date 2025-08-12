package storage

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

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
		return NewGormDB(dsn)

	case "pq":
		return NewPqDB(dsn)

	default:
		return nil, fmt.Errorf("unsupported DB engine: %s", dbEngine)
	}
}
