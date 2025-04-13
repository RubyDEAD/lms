package main

//sudo docker run -d --hostname rabbitlms --name rabbitlms -p 15672:15672 -p 5672:5672 rabbitmq:3-management
import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
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

	// Declare the request queue (same as client)
	_, err = ch.QueueDeclare("testReq", false, false, false, false, nil)
	if err != nil {
		log.Fatal("Failed to declare queue:", err)
	}

	msgs, err := ch.Consume(
		"testReq", "", true, false, false, false, nil)
	if err != nil {
		log.Fatal("Failed to register consumer:", err)
	}

	// Set up graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Waiting for requests...")

	for {
		select {
		case msg := <-msgs:
			response := fmt.Sprintf("Processed: %s (Correlation ID: %s)",
				msg.Body, msg.CorrelationId)

			err := ch.Publish(
				"", msg.ReplyTo, false, false,
				amqp.Publishing{
					Body:          []byte(response),
					CorrelationId: msg.CorrelationId,
				},
			)
			if err != nil {
				log.Printf("Failed to send response: %v", err)
			}

			fmt.Printf("Request received: %s\n", msg.Body)

		case <-sigChan:
			log.Println("Shutting down server...")
			return
		}
	}
}
