package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"borrowing-service/graph"
	"borrowing-service/graph/generated"
	"borrowing-service/internal/bookservice"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

const (
	defaultPort    = "8080"
	maxDBRetries   = 5
	retryDelay     = 5 * time.Second
)

func main() {
	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Connect to borrowing-service's Supabase
	borrowingDB := connectToSupabase("SUPABASE_DB_URL", logger)
	defer borrowingDB.Close()

	// Initialize book-service client (using its Supabase URL)
	bookServiceClient := bookservice.NewClient(bookservice.Config{
		BaseURL: os.Getenv("BOOK_SERVICE_SUPABASE_URL"), // Different env var
		Timeout: 5 * time.Second,
		Logger:  logger,
	})

	// Create resolver with dependencies
	resolver := &graph.Resolver{
		DB:          borrowingDB,
		BookService: bookServiceClient,
		Logger:      logger,
	}

	// GraphQL server setup
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: resolver,
	}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	logger.Info("Server starting", zap.String("port", port))
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func connectToSupabase(envVar string, logger *zap.Logger) *pgxpool.Pool {
	dbURL := os.Getenv(envVar)
	if dbURL == "" {
		logger.Fatal("Supabase URL not set", zap.String("env_var", envVar))
	}

	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		logger.Fatal("Failed to parse database config", zap.Error(err))
	}

	// Connection pool settings
	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	var dbpool *pgxpool.Pool

	// Retry logic
	for i := 0; i < maxDBRetries; i++ {
		dbpool, err = pgxpool.NewWithConfig(context.Background(), config)
		if err == nil {
			break
		}

		logger.Warn("Connection attempt failed",
			zap.Int("attempt", i+1),
			zap.Error(err))
		
		if i < maxDBRetries-1 {
			time.Sleep(retryDelay * time.Duration(i+1))
		}
	}

	if err != nil {
		logger.Fatal("Failed to connect after retries", zap.Error(err))
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := dbpool.Ping(ctx); err != nil {
		logger.Fatal("Failed to ping database", zap.Error(err))
	}

	logger.Info("Successfully connected to Supabase")
	return dbpool
}