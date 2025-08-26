package storage

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func InitRedis() *redis.Client {
	godotenv.Load()
	dbNum, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       dbNum,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("❌ Failed to connect Redis: %v", err)
	}

	log.Println("✅ Redis connected")
	return client
}
