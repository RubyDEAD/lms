package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/RubyDEAD/lms/borrowing-service/graph"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5"
	"github.com/nedpals/supabase-go"
	"github.com/rs/cors"
	"github.com/vektah/gqlparser/v2/ast"
	// Import the Supabase Go client library
)

const defaultPort = "8082"

func main() {
	//dbURL := "postgresql://postgres:FGar.Uzebyg3ZZ9@db.ictfypsqogdoceosoqdj.supabase.co:5432/postgres"
	dbURL := os.Getenv("DATABASE_URL")
	// Create a database connection pool
	db, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer db.Close(context.Background())

	fmt.Println("Successfully connected to Supabase using pgx!")

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	supabaseURL := "https://db.ictfypsqogdoceosoqdj.supabase.co"
	supabaseKey := os.Getenv("SUPABASE_KEY")
	if supabaseKey == "" {
		log.Fatal("SUPABASE_KEY environment variable is not set")
	}

	// Create Supabase client
	supabaseClient := supabase.CreateClient(supabaseURL, supabaseKey)

	resolver := graph.NewResolver(db, supabaseClient)
	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	srv.AddTransport(&transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		KeepAlivePingInterval: 10 * time.Second,
	})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	// Enable CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8081"}, // Allow requests from your front-end
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"}, // Allow specific HTTP methods
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
	}).Handler(srv)

	// Use corsHandler for the /query endpoint
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", corsHandler)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
