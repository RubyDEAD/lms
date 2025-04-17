package patronmq

//sudo docker run -d --hostname rabbitlms --name rabbitlms -p 15672:15672 -p 5672:5672 rabbitmq:3-management
import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/GSalise/lms/patron-service/graph"
	"github.com/jackc/pgx/v5/pgxpool"
	amqp "github.com/rabbitmq/amqp091-go"
)

type GraphQLMessage struct {
	Query             string                 `json:"query"`
	Variables         map[string]interface{} `json:"variables"`
	RequestedResolver string                 `json:"requestedResolver"`
}

func StartRabbitMQConsumer(dbpool *pgxpool.Pool) {
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

	_, err = ch.QueueDeclare(
		"patron-service-queue",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatal("Failed to declare queue:", err)
	}

	msgs, err := ch.Consume(
		"patron-service-queue",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
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
			var data GraphQLMessage
			if err := json.Unmarshal(msg.Body, &data); err != nil {
				log.Printf("Error decoding JSON: %v", err)
				continue // Skip malformed messages
			}

			resolver := &graph.Resolver{
				DB: dbpool,
			}
			ctx := context.Background()

			switch data.RequestedResolver {
			case "createPatron":
				firstName, _ := data.Variables["firstName"].(string)
				lastName, _ := data.Variables["lastName"].(string)
				phoneNumber, _ := data.Variables["phoneNumber"].(string)

				patron, ResolverErr := resolver.Mutation().CreatePatron(ctx, firstName, lastName, phoneNumber)

				if ResolverErr != nil {
					log.Fatalf("err: %v", ResolverErr)
				}

				response := map[string]interface{}{
					"data": map[string]interface{}{
						"createPatron": patron,
					},
				}

				result, err := json.Marshal(response)
				if err != nil {
					log.Fatalf("Error marshalling patron to JSON: %v", err)
				}
				log.Printf("reply to : %s", msg.ReplyTo)
				log.Printf("reply to : %s", result)
				log.Printf("reply to : %s", msg.CorrelationId)
				err = ch.Publish(
					"",
					msg.ReplyTo,
					false,
					false,
					amqp.Publishing{
						ContentType:   "application/json",
						CorrelationId: msg.CorrelationId,
						Body:          result,
					},
				)

				if err != nil {
					log.Fatalf("failed to publish message: %v", err)
				}

			default:
				log.Printf("Unknown resolver")
			}

		case <-sigChan:
			log.Println("Shutting down server...")
			return
		}
	}

}
