package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Declare the request queue this wer you send the availability request
	_, err = ch.QueueDeclare(
		"bookCopyAvailRequests", // name
		true,                    // durable
		false,                   // delete when unused
		false,                   // exclusive
		false,                   // no-wait
		nil,                     // arguments
	)
	failOnError(err, "Failed to declare the request queue")

	// Declare the reply queue this were you receive the replies for the availability request
	replyQueue, err := ch.QueueDeclare(
		"bookServiceReplies", // name
		true,                 // durable
		false,                // delete when unused
		false,                // exclusive
		false,                // no-wait
		nil,                  // arguments
	)
	failOnError(err, "Failed to declare the reply queue")

	// Send a borrow request
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create a JSON payload for the borrow request
	borrowRequest := struct {
		BookID string `json:"book_id"`
	}{
		BookID: "01bb9f98-4c01-448b-84aa-bf01f2c2b3e0", // Example book ID
	}

	requestBody, err := json.Marshal(borrowRequest)
	failOnError(err, "Failed to marshal borrow request payload")

	err = ch.PublishWithContext(ctx,
		"",                      // exchange
		"bookCopyAvailRequests", // routing key
		false,                   // mandatory
		false,                   // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        requestBody,
		})
	failOnError(err, "Failed to publish borrow request")
	log.Printf(" [x] Sent borrow request: %s\n", requestBody)

	// Consume the reply . the bookService sends back the availabile book copies
	msgs, err := ch.Consume(
		replyQueue.Name, // queue name
		"",              // consumer tag
		true,            // auto-ack
		false,           // exclusive
		false,           // no-local
		false,           // no-wait
		nil,             // arguments
	)
	failOnError(err, "Failed to register a consumer for replies")

	// Listen for replies
	for d := range msgs {
		log.Printf("Received a reply: %s", d.Body)

		// Parse the reply
		var reply struct {
			BookCopyID string `json:"book_copy_id"`
			Status     string `json:"status"`
		}
		if err := json.Unmarshal(d.Body, &reply); err != nil {
			log.Printf("Failed to parse reply: %v", err)
			continue
		}

		// Check if the book copy is available
		if reply.Status == "Available" {
			log.Printf("Book copy is available. Requesting status update...")

			// Send a request to update the book copy status
			updatePayload := struct {
				ID     string  `json:"id"`
				Status *string `json:"status"`
			}{
				ID:     reply.BookCopyID,
				Status: stringPtr("Borrowed"), // Update status to "Borrowed"
			}

			updateBody, err := json.Marshal(updatePayload)
			failOnError(err, "Failed to marshal update payload")

			err = ch.PublishWithContext(ctx,
				"",                  // exchange
				"book-copies-queue", // routing key
				false,               // mandatory
				false,               // immediate
				amqp.Publishing{
					ContentType: "application/json",
					Body:        updateBody,
				})
			failOnError(err, "Failed to publish update request")
			log.Printf(" [x] Sent update request: %s\n", updateBody)
		} else {
			log.Printf("no Available copies")
		}
		break // Exit after processing one reply
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func stringPtr(s string) *string {
	return &s
}
