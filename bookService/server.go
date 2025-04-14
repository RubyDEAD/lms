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
	"github.com/Cat6utpcableclarke/bookService/graph"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8080"

func main() {

	dbURL := "postgresql://postgres.hwkuzfsecehszlftxqpn:cat6utpcable@aws-0-ap-southeast-1.pooler.supabase.com:5432/postgres"

	// Create a database connection pool
	db, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatalf(" Unable to connect to database: %v", err)
	}
	defer db.Close(context.Background())

	fmt.Println("Successfully connected to Supabase using pgx!")

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	resolver := graph.NewResolver(db)
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

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
