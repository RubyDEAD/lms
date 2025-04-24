package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Cat6utpcableclarke/API-gateway/graph"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8081"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Create the GraphQL server
	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	// Add WebSocket transport for subscriptions
	srv.AddTransport(&transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins (adjust for production)
			},
		},
		KeepAlivePingInterval: 10 * time.Second,
	})

	// Add HTTP transport for queries and mutations
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.Options{})

	// Set query caching
	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	// Enable introspection and persisted queries
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	// Enable CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Allow requests from your front-end
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"}, // Allow specific HTTP methods
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
	}).Handler(srv)

	// Set up HTTP handlers
	http.Handle("/", playground.Handler("GraphQL Playground", "/query"))
	http.Handle("/query", corsHandler)

	// Start the server
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
