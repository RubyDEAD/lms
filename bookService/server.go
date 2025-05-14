package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	consumer "github.com/Cat6utpcableclarke/bookService/Consumer"
	db "github.com/Cat6utpcableclarke/bookService/db"
	"github.com/Cat6utpcableclarke/bookService/graph"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8080"

func main() {

	// Create a database connection pool
	db, err := db.ConnectPool()
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	conn, err := db.Acquire(context.Background())
	if err != nil {
		log.Fatalf("Unable to acquire a database connection: %v", err)
	}
	defer conn.Release()

	resolver := graph.NewResolver(db)
	go func() {
		log.Println("Listening for book copy updates")
		consumer.UpdateConsumer(resolver)
	}()

	go func() {
		log.Println("Listening for availability requests")
		consumer.ListenAvailabiltyRequests(resolver)
	}()
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
