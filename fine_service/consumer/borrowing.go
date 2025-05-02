package consumer

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type BorrowingReturnedEvent struct {
	BorrowingID string `json:"borrowingId"`
	PatronID    string `json:"patronId"`
	BookID      string `json:"bookId"`
	DueDate     string `json:"dueDate"`     // ISO8601 format (e.g., "2025-04-28")
	ReturnedAt  string `json:"returnedAt"`  // ISO8601 format (e.g., "2025-05-01")
}

// ListenBorrowingReturned starts a consumer that listens for borrowing.returned events
func ListenBorrowingReturned(ch *amqp.Channel, db *sql.DB) {
	msgs, err := ch.Consume(
		"borrowing.returned", // queue name
		"",                   // consumer
		true,                 // auto-ack
		false,                // exclusive
		false,                // no-local
		false,                // no-wait
		nil,                  // args
	)
	if err != nil {
		log.Fatalf("Failed to register RabbitMQ consumer: %v", err)
	}

	go func() {
		for d := range msgs {
			var event BorrowingReturnedEvent
			err := json.Unmarshal(d.Body, &event)
			if err != nil {
				log.Printf("Failed to parse event: %v", err)
				continue
			}

			log.Printf("Received borrowing.returned event: %+v", event)

			// Parse dates
			dueDate, _ := time.Parse(time.RFC3339, event.DueDate)
			returnedAt, _ := time.Parse(time.RFC3339, event.ReturnedAt)

			if returnedAt.After(dueDate) {
				daysLate := int(returnedAt.Sub(dueDate).Hours() / 24)
				amount := float64(daysLate) * 1.0 // e.g., $1 per day

				// Insert fine into DB
				_, err := db.ExecContext(context.Background(),
					`INSERT INTO fines (borrowing_id, patron_id, book_id, days_late, amount, created_at)
					 VALUES ($1, $2, $3, $4, $5, NOW())`,
					event.BorrowingID, event.PatronID, event.BookID, daysLate, amount,
				)
				if err != nil {
					log.Printf("Failed to insert fine: %v", err)
				} else {
					log.Printf("Fine created for borrowing %s: $%.2f", event.BorrowingID, amount)
				}
			} else {
				log.Printf("Book returned on time. No fine.")
			}
		}
	}()

	log.Println("Listening for borrowing.returned events...")
}
