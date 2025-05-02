package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"fine_service/graph"
	"fine_service/consumer"
	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8000"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Setup DB connection
	db, err := sql.Open("postgres", "postgresql://postgres.tltusctrslkkwukzfoib:NXH3QMNg3IGSpBAZ@aws-0-ap-southeast-1.pooler.supabase.com:5432/postgres")
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer db.Close()

	// ✅ Connect to RabbitMQ
	rabbitConn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitConn.Close()

	rabbitCh, err := rabbitConn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer rabbitCh.Close()

	// // ✅ Declare queue (if not already declared by borrowing-service)
	// _, err = rabbitCh.QueueDeclare(
	// 	"borrowing.returned", // queue name
	// 	true,                 // durable
	// 	false,                // auto-delete
	// 	false,                // exclusive
	// 	false,                // no-wait
	// 	nil,                  // arguments
	// )
	// if err != nil {
	// 	log.Fatalf("Failed to declare queue: %v", err)
	// }

	// // Declare a queue (optional)
	// _, err = rabbitCh.QueueDeclare(
	// 	"fine_created_queue", // name
	// 	true,                 // durable
	// 	false,                // delete when unused
	// 	false,                // exclusive
	// 	false,                // no-wait
	// 	nil,                  // arguments
	// )
	// if err != nil {
	// 	log.Fatalf("Failed to declare queue: %v", err)
	// }

	// ✅ Start listening for borrowing.returned events
	// go consumer.FineService{DB: db}.ListenBorrowingReturned(rabbitConn)

	consumer.ListenBorrowingReturned(rabbitCh, db)

	// ✅ Initialize GraphQL server
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{
			DB:            db,
			RabbitChannel: rabbitCh, // ✅ Make sure Resolver struct has this
		},
	}))

	srv.AddTransport(&transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	go consumer.ListenBorrowingReturned(rabbitCh, db)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}


