package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	// "os/signal"
	// "sync"
	// "syscall"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/GSalise/lms/patron-service/graph"
	"github.com/GSalise/lms/patron-service/graph/model"
	"github.com/rs/cors"

	"github.com/GSalise/lms/patron-service/patronmq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8420"

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Supabase database connection string
	dbURL := os.Getenv("SUPABASE_DB_URL")
	if dbURL == "" {
		log.Fatal("SUPABASE_DB_URL environment variable is not set")
	}

	// Initialize the database connection pool
	dbpool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer dbpool.Close()

	// Pass the database connection to the Resolver
	resolver := &graph.Resolver{
		DB:                dbpool,
		PatronSubscribers: make(map[chan *model.Patron]bool),
	}

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
	})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8069"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
	}).Handler(srv)

	// http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	// http.Handle("/query", srv)

	// log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	// log.Fatal(http.ListenAndServe(":"+port, nil))

	// Use a WaitGroup to manage concurrency
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the GraphQL server in a goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()

		server := &http.Server{Addr: ":" + port}

		http.Handle("/", playground.Handler("GraphQL playground", "/query"))
		http.Handle("/query", corsHandler)

		go func() {
			log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("GraphQL server failed: %v", err)
			}
		}()

		// Wait for context cancellation to shut down the server
		<-ctx.Done()
		log.Println("Shutting down GraphQL server...")
		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("GraphQL server shutdown error: %v", err)
		}
	}()

	// Start the message queue in a goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		patronmq.StartRabbitMQConsumer(dbpool)
	}()

	// Listen for OS signals (e.g., SIGINT, SIGTERM)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan // Blocks until a signal is received
	cancel()  // Notify goroutines to clean up
	log.Println("Shutting down...")

	wg.Wait() // Wait for all goroutines to finish
	log.Println("Server gracefully stopped.")
}
