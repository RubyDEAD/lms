package server

import (
	"book-service/gql"
	"book-service/repository"
	"context"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/jackc/pgx/v5"
)

func StartServer(dbConn *pgx.Conn) {
	repo := repository.NewBookRepository(dbConn)
	resolver := gql.Resolver{Repo: repo}

	srv := handler.NewDefaultServer(gql.NewExecutableSchema(gql.Config{Resolvers: &resolver}))

	http.Handle("/", playground.Handler("GraphQL Playground", "/query"))
	http.Handle("/query", srv)

	log.Println("🚀 Server running on http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
