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
	"github.com/gorilla/websocket"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/rs/cors"
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

	// Setup GraphQL server with custom transports for subscriptions
srv := handler.New(graph.NewExecutableSchema(graph.Config{
	Resolvers: &graph.Resolver{
		DB:            db,
		RabbitChannel: rabbitCh,
	},
}))

srv.AddTransport(transport.Websocket{
	Upgrader: websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// Allow all origins for local testing (secure this in production)
			return true
		},
	},
	KeepAlivePingInterval: 10 * time.Second,
})

srv.AddTransport(transport.Options{})
srv.AddTransport(transport.GET{})
srv.AddTransport(transport.POST{})
srv.AddTransport(transport.MultipartForm{})

http.Handle("/", playground.Handler("GraphQL playground", "/query"))
http.Handle("/query", srv)


		c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // frontend origin
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
	})

	server := &http.Server{
		Addr:    ":" + port,
		Handler: c.Handler(http.DefaultServeMux), // <- wrap default mux
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




