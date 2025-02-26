package main

import (
	"book-service/graph" // Import the graph package, which includes generated.go
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQL driver (used implicitly by database/sql)
)

const defaultPort = "8080"

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using default environment variables")
	}

	// Get the port from environment variables or use the default
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Get the database URL from environment variables
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set in .env")
	}

	// Connect to the database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Ping the database to verify the connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	fmt.Println("Successfully connected to Supabase!")

	// Create a new GraphQL server
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	// Add supported transports
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	// Set up query caching with LRU cache
	srv.SetQueryCache(lru.New[*graph.Response](1000)) // Explicitly specify the type parameter

	// Enable introspection and APQ (Automatic Persisted Queries)
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{Cache: lru.New[*graph.Response](100)}) // Explicitly specify the type parameter

	// Set up the GraphQL Playground
	http.Handle("/", playground.Handler("GraphQL Playground", "/query"))
	http.Handle("/query", srv)

	// Start the server
	log.Printf("Connect to http://localhost:%s/ for GraphQL Playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}