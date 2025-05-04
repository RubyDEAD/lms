// package main

// import (
// 	"database/sql"
// 	"log"
// 	"net/http"
// 	"os"
// 	"time"

// 	"fine_service/graph"
// 	"fine_service/consumer"
// 	_ "github.com/lib/pq"
// 	amqp "github.com/rabbitmq/amqp091-go"
// 	"github.com/99designs/gqlgen/graphql/handler"
// 	"github.com/99designs/gqlgen/graphql/handler/transport"
// 	"github.com/99designs/gqlgen/graphql/playground"
// )

// const defaultPort = "8000"

// func main() {
// 	port := os.Getenv("PORT")
// 	if port == "" {
// 		port = defaultPort
// 	}

// 	// Setup DB connection
// 	db, err := sql.Open("postgres", "postgresql://postgres.tltusctrslkkwukzfoib:NXH3QMNg3IGSpBAZ@aws-0-ap-southeast-1.pooler.supabase.com:5432/postgres")
// 	if err != nil {
// 		log.Fatalf("failed to connect to db: %v", err)
// 	}
// 	defer db.Close()

// 	// ✅ Connect to RabbitMQ
// 	rabbitConn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
// 	if err != nil {
// 		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
// 	}
// 	defer rabbitConn.Close()

// 	rabbitCh, err := rabbitConn.Channel()
// 	if err != nil {
// 		log.Fatalf("Failed to open a channel: %v", err)
// 	}
// 	defer rabbitCh.Close()

// 	// ✅ Declare the "borrowing.returned" queue
// 	_, err = rabbitCh.QueueDeclare(
// 		"borrowing.returned", // queue name
// 		true,                 // durable
// 		false,                // auto-delete
// 		false,                // exclusive
// 		false,                // no-wait
// 		nil,                  // arguments
// 	)
// 	if err != nil {
// 		log.Fatalf("Failed to declare queue: %v", err)
// 	}

// 	// ✅ Start listening for borrowing.returned events
// 	go consumer.ListenBorrowingReturned(rabbitCh, db)

// 	// ✅ Initialize GraphQL server
// 	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
// 		Resolvers: &graph.Resolver{
// 			DB:            db,
// 			RabbitChannel: rabbitCh, // ✅ Make sure Resolver struct has this
// 		},
// 	}))

// 	srv.AddTransport(&transport.Websocket{
// 		KeepAlivePingInterval: 10 * time.Second,
// 	})

// 	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
// 	http.Handle("/query", srv)

// 	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
// 	log.Fatal(http.ListenAndServe(":"+port, nil))
// }

package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fine_service/consumer"
	"fine_service/graph"

	_ "github.com/lib/pq"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	amqp "github.com/rabbitmq/amqp091-go"
)

const defaultPort = "8000"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Connect to PostgreSQL
	db, err := sql.Open("postgres", "postgresql://postgres.tltusctrslkkwukzfoib:NXH3QMNg3IGSpBAZ@aws-0-ap-southeast-1.pooler.supabase.com:5432/postgres")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Connect to RabbitMQ
	rabbitConn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitConn.Close()

	rabbitCh, err := rabbitConn.Channel()
	if err != nil {
		log.Fatalf("Failed to open RabbitMQ channel: %v", err)
	}
	defer rabbitCh.Close()

	// Start the return event consumer
	go func() {
		err := consumer.ConsumeReturnEvent()
		if err != nil {
			log.Fatalf("Error consuming return events: %v", err)
		}
	}()

	// Setup GraphQL server
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{
			DB:            db,
			RabbitChannel: rabbitCh,
		},
	}))

	srv.AddTransport(&transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	// Graceful shutdown
	server := &http.Server{
		Addr:    ":" + port,
		Handler: nil,
	}

	go func() {
		log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Listen for interrupt signal to gracefully shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	log.Println("Shutting down fine-service...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP server Shutdown: %v", err)
	}

	log.Println("fine-service shutdown complete.")
}



