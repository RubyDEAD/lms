package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func CheckAvailability(BookID string) (bool, *amqp.Connection, string, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return false, nil, "", fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return false, nil, "", fmt.Errorf("failed to open a channel: %v", err)
	}
	defer ch.Close()

	// Declare the request queue
	_, err = ch.QueueDeclare(
		"bookCopyAvailRequests", // name
		true,                    // durable
		false,                   // delete when unused
		false,                   // exclusive
		false,                   // no-wait
		nil,                     // arguments
	)
	if err != nil {
		conn.Close()
		return false, nil, "", fmt.Errorf("failed to declare the request queue: %v", err)
	}

	// Declare the reply queue
	replyQueue, err := ch.QueueDeclare(
		"bookServiceReplies", // name
		true,                 // durable
		false,                // delete when unused
		false,                // exclusive
		false,                // no-wait
		nil,                  // arguments
	)
	if err != nil {
		conn.Close()
		return false, nil, "", fmt.Errorf("failed to declare the reply queue: %v", err)
	}

	// Send a borrow request
	err = sendAvailabilityRequest(ch, BookID)
	if err != nil {
		conn.Close()
		return false, nil, "", fmt.Errorf("failed to send availability request: %v", err)
	}

	// Consume the reply
	msgs, err := ch.Consume(
		replyQueue.Name, // queue name
		"",              // consumer tag
		true,            // auto-ack
		false,           // exclusive
		false,           // no-local
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		conn.Close()
		return false, nil, "", fmt.Errorf("failed to register a consumer for replies: %v", err)
	}

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
			log.Printf("Book copy is available. Returning BookCopyID: %s", reply.BookCopyID)
			return true, conn, reply.BookCopyID, nil
		} else {
			log.Printf("No available book copies. Status: %s", reply.Status)
			conn.Close()
			return false, nil, "", nil
		}
	}

	conn.Close()
	return false, nil, "", nil
}

// sendAvailabilityRequest sends a borrow request to check for available book copies
func sendAvailabilityRequest(ch *amqp.Channel, bookID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create a JSON payload for the borrow request
	borrowRequest := struct {
		BookID string `json:"book_id"`
	}{
		BookID: bookID,
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
	return nil
}

// sendUpdateRequest sends a request to update the status of a book copy
func SendUpdateRequest(conn *amqp.Connection, bookCopyID string, status string) error {
	// Open a new channel
	ch, err := conn.Channel()
	if err != nil {
		log.Printf("Failed to open a channel: %v", err)
		return fmt.Errorf("failed to open a channel: %v", err)
	}
	defer ch.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create a JSON payload for the update request
	updatePayload := struct {
		ID     string  `json:"id"`
		Status *string `json:"status"`
	}{
		ID:     bookCopyID,
		Status: stringPtr(status),
	}

	updateBody, err := json.Marshal(updatePayload)
	if err != nil {
		log.Printf("Failed to marshal update payload: %v", err)
		return fmt.Errorf("failed to marshal update payload: %v", err)
	}

	// Publish the update request to RabbitMQ
	log.Printf("Publishing update request for BookCopyID: %s with status: %s", bookCopyID, status)
	err = ch.PublishWithContext(ctx,
		"",                  // exchange
		"book-copies-queue", // routing key
		false,               // mandatory
		false,               // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        updateBody,
		})
	if err != nil {
		log.Printf("Failed to publish update request: %v", err)
		return fmt.Errorf("failed to publish update request: %v", err)
	}

	log.Printf(" [x] Successfully sent update request: %s\n", updateBody)
	return nil
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Printf("%s: %s", msg, err)
	}
}

func stringPtr(s string) *string {
	return &s
}
