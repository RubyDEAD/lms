package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// Connect to RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open channel:", err)
	}
	defer ch.Close()

	// Declare queues
	RequestQueue, err := ch.QueueDeclare(
		"testReq", false, false, false, false, nil)
	if err != nil {
		log.Fatal("Failed to declare request queue:", err)
	}

	ResponseQueue, err := ch.QueueDeclare(
		"testResp", false, false, false, false, nil)
	if err != nil {
		log.Fatal("Failed to declare response queue:", err)
	}

	// Start consuming responses FIRST
	msgs, err := ch.Consume(
		ResponseQueue.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal("Failed to register consumer:", err)
	}

	// Set up graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Send 5 requests
	for i := 0; i < 5; i++ {
		corrID := uuid.New().String()
		msg := fmt.Sprintf("hello world %d", i)

		err = ch.Publish(
			"", RequestQueue.Name, false, false,
			amqp.Publishing{
				ContentType:   "text/plain",
				Body:          []byte(msg),
				ReplyTo:       ResponseQueue.Name,
				CorrelationId: corrID,
			},
		)
		if err != nil {
			log.Printf("Failed to publish message %d: %v", i, err)
			continue
		}
		log.Printf("Sent request %d with correlation ID: %s", i, corrID)
	}

	// Process responses with timeout
	timeout := time.After(10 * time.Second)

	for {
		select {
		case msg := <-msgs:
			log.Printf("Received response: %s (Correlation ID: %s)",
				msg.Body, msg.CorrelationId)

		case <-timeout:
			log.Println("Timeout waiting for responses")
			return

		case <-sigChan:
			log.Println("Shutting down...")
			return
		}
	}
}
