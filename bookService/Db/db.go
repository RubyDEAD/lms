package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func ConnectPool() (*pgxpool.Pool, error) {
	// Load environment variables from .env file
	err := godotenv.Load("DB/.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get the DATABASE_URL from the environment
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatalf("DATABASE_URL environment variable is not set")
	}

	// Configure the connection pool
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatalf("Unable to parse database URL: %v", err)
		return nil, err
	}

	// Set connection pool settings
	config.MaxConns = 10                 // Maximum number of connections
	config.MinConns = 2                  // Minimum number of connections
	config.MaxConnIdleTime = time.Minute // Maximum idle time for connections

	// Create the connection pool
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
		return nil, err
	}

	fmt.Println("Successfully connected to Supabase using pgxpool!")
	return pool, nil
}
