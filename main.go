package main

import (
	"context"
	"log"
	"os"

	_ "github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	_ "github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/require"
	"smartrb.com/api"
	"smartrb.com/db"
	"smartrb.com/mq"
	_ "smartrb.com/users"
	_ "smartrb.com/categories"
	_ "smartrb.com/interests"
	_ "smartrb.com/flashcards"
	_ "smartrb.com/fill_in_blanks"
)

func main() {
	// Load .env file if it exists, but don't fail if it doesn't
	// This allows local development with .env and production with environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Note: No .env file found, using environment variables:", err)
		// Don't call log.Fatal() - continue with environment variables
	}

	ctx := context.Background()

	// Setup DB instance
	sqlDB, db_err := db.New(ctx, os.Getenv("DB_CONNECTION_STRING"))
	if db_err != nil {
		log.Fatal("Database connection error:", db_err)
	}
	defer sqlDB.Close()

	db.SetDefault(sqlDB)

	// Setup MQ instance
	amqpURL := os.Getenv("AMQP_URL") // e.g., amqp://guest:guest@localhost:5672/

	mqClient, mq_err := mq.NewClient(ctx, amqpURL)
	if mq_err != nil {
		log.Fatal("RabbitMQ connection error:", mq_err)
	}
	// defer mqClient.Close()
	mq.SetDefault(mqClient)

	// Start API server
	api.StartAPI()
}