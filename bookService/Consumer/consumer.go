package consumer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Cat6utpcableclarke/bookService/graph"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	rabbitMQURL     = "amqp://guest:guest@localhost:5672/"
	bookCopiesQueue = "book-copies-queue"
)

func InitConsumer(resolver *graph.Resolver) {
	conn, err := amqp.Dial(rabbitMQURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	_, err = ch.QueueDeclare(
		bookCopiesQueue, // queue name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		bookCopiesQueue, // queue name
		"",              // consumer tag
		true,            // auto-ack
		false,           // exclusive
		false,           // no-local
		false,           // no-wait
		nil,             // arguments
	)
	failOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			var payload struct {
				ID     string  `json:"id"`
				Status *string `json:"status"`
			}
			if err := json.Unmarshal(d.Body, &payload); err != nil {
				log.Printf("Failed to parse message: %v", err)
				continue
			}

			//Call the UpdateBookCopyStatus resolver
			ctx := context.Background()
			updatedCopy, err := resolver.Mutation().UpdateBookCopyStatus(ctx, payload.ID, payload.Status)
			if err != nil {
				log.Printf("Failed to update book copy status: %v", err)
				continue
			}

			log.Printf("Successfully updated book copy: %+v", updatedCopy)
		}
	}()

	log.Println("Waiting for messages. To exit press CTRL+C")
	select {}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
