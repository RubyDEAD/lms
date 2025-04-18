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
	"github.com/GSalise/lms/patron-service/graph/model"
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
					log.Fatalf("Create Patron Resolver err: %v", ResolverErr)
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
				// log.Printf("reply to : %s", msg.ReplyTo)
				// log.Printf("reply to : %s", result)
				// log.Printf("reply to : %s", msg.CorrelationId)
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
			case "updatePatron":
				patron_id, _ := data.Variables["patron_id"].(string)
				firstName, _ := data.Variables["firstName"].(string)
				lastName, _ := data.Variables["lastName"].(string)
				phoneNumber, _ := data.Variables["phoneNumber"].(string)

				patron, ResolverErr := resolver.Mutation().UpdatePatron(ctx, patron_id, &firstName, &lastName, &phoneNumber)

				if ResolverErr != nil {
					log.Fatalf("Update Patron Resolver err: %v", ResolverErr)
				}

				response := map[string]interface{}{
					"data": map[string]interface{}{
						"updatePatron": patron,
					},
				}

				result, err := json.Marshal(response)
				if err != nil {
					log.Fatalf("Error marshalling patron to JSON: %v", err)
				}
				// log.Printf("reply to : %s", msg.ReplyTo)
				// log.Printf("reply to : %s", result)
				// log.Printf("reply to : %s", msg.CorrelationId)
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

			case "deletePatronById":
				patron_id, _ := data.Variables["patron_id"].(string)

				patron, ResolverErr := resolver.Mutation().DeletePatronByID(ctx, patron_id)

				if ResolverErr != nil {
					log.Fatalf("err: %v", ResolverErr)
				}

				response := map[string]interface{}{
					"data": map[string]interface{}{
						"deletePatronById": patron,
					},
				}

				result, err := json.Marshal(response)
				if err != nil {
					log.Fatalf("Error marshalling patron to JSON: %v", err)
				}

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

			case "updateMembershipByPatronId":
				patron_id, _ := data.Variables["patron_id"].(string)
				level, _ := data.Variables["level"].(string)

				patron, ResolverErr := resolver.Mutation().UpdateMembershipByPatronID(ctx, patron_id, model.MembershipLevel(level))

				if ResolverErr != nil {
					log.Fatalf("Update Membershup By Patron ID Resolver err: %v", ResolverErr)
				}

				response := map[string]interface{}{
					"data": map[string]interface{}{
						"updateMembershipByPatronId": patron,
					},
				}

				result, err := json.Marshal(response)
				if err != nil {
					log.Fatalf("Error marshalling patron to JSON: %v", err)
				}

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

			case "updateMembershipByMembershipId":
				membership_id, _ := data.Variables["membership_id"].(string)
				level, _ := data.Variables["level"].(string)

				patron, ResolverErr := resolver.Mutation().UpdateMembershipByMembershipID(ctx, membership_id, model.MembershipLevel(level))

				if ResolverErr != nil {
					log.Fatalf("Update Membership By Membership ID Resolver err: %v", ResolverErr)
				}

				response := map[string]interface{}{
					"data": map[string]interface{}{
						"updateMembershipByMembershipId": patron,
					},
				}

				result, err := json.Marshal(response)
				if err != nil {
					log.Fatalf("Error marshalling patron to JSON: %v", err)
				}

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

			case "updatePatronStatus":
				patron_id, _ := data.Variables["patron_id"].(string)
				var warning_count *int32
				if raw, ok := data.Variables["warning_count"]; ok && raw != nil {
					if val, ok := raw.(float64); ok {
						conv := int32(val)
						warning_count = &conv
					}
				}

				var unpaid_fees *float64
				if raw, ok := data.Variables["unpaid_fees"]; ok && raw != nil {
					if val, ok := raw.(float64); ok {
						unpaid_fees = &val
					}
				}

				var status *string
				if raw, ok := data.Variables["status"].(string); ok && raw != "" {
					status = &raw
				}

				patron, ResolverErr := resolver.Mutation().UpdatePatronStatus(ctx, patron_id, warning_count, unpaid_fees, (*model.Status)(status))

				if ResolverErr != nil {
					log.Fatalf("Update Patron Status Resolver err: %v", ResolverErr)
				}

				response := map[string]interface{}{
					"data": map[string]interface{}{
						"updatePatronStatus": patron,
					},
				}

				result, err := json.Marshal(response)
				if err != nil {
					log.Fatalf("Error marshalling patron to JSON: %v", err)
				}

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

			case "addViolation":
				patron_id, _ := data.Variables["patron_id"].(string)
				violation_type, _ := data.Variables["violation_type"].(string)
				violation_info, _ := data.Variables["violation_info"].(string)

				patron, ResolverErr := resolver.Mutation().AddViolation(ctx, patron_id, model.ViolationType(violation_type), violation_info)

				if ResolverErr != nil {
					log.Fatalf("Add Violation Resolver err: %v", ResolverErr)
				}

				response := map[string]interface{}{
					"data": map[string]interface{}{
						"addViolation": patron,
					},
				}

				result, err := json.Marshal(response)
				if err != nil {
					log.Fatalf("Error marshalling patron to JSON: %v", err)
				}

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

			case "updateViolationStatus":
				violation_id, _ := data.Variables["violation_id"].(string)
				violation_status, _ := data.Variables["violation_type"].(string)

				patron, ResolverErr := resolver.Mutation().UpdateViolationStatus(ctx, violation_id, model.ViolationStatus(violation_status))

				if ResolverErr != nil {
					log.Fatalf("Update Violation Status Resolver err: %v", ResolverErr)
				}

				response := map[string]interface{}{
					"data": map[string]interface{}{
						"updateViolationStatus": patron,
					},
				}

				result, err := json.Marshal(response)
				if err != nil {
					log.Fatalf("Error marshalling patron to JSON: %v", err)
				}

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

			case "getPatronById":
				patron_id, _ := data.Variables["patron_id"].(string)

				patron, ResolverErr := resolver.Query().GetPatronByID(ctx, patron_id)
				if ResolverErr != nil {
					log.Fatalf("Get Patron By ID Resolver err: %v", ResolverErr)
				}

				response := map[string]interface{}{
					"data": map[string]interface{}{
						"getPatronById": patron,
					},
				}

				result, err := json.Marshal(response)
				if err != nil {
					log.Fatalf("Error marshalling patron to JSON: %v", err)
				}

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

			case "getAllPatrons":
				patron, ResolverErr := resolver.Query().GetAllPatrons(ctx)
				if ResolverErr != nil {
					log.Fatalf("Get Patron By ID Resolver err: %v", ResolverErr)
				}

				response := map[string]interface{}{
					"data": map[string]interface{}{
						"getAllPatrons": patron,
					},
				}

				result, err := json.Marshal(response)
				if err != nil {
					log.Fatalf("Error marshalling patron to JSON: %v", err)
				}

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

			case "getMembershipByLevel":
				level, _ := data.Variables["level"].(string)

				patron, ResolverErr := resolver.Query().GetMembershipByLevel(ctx, model.MembershipLevel(level))
				if ResolverErr != nil {
					log.Fatalf("Get Patron By ID Resolver err: %v", ResolverErr)
				}

				response := map[string]interface{}{
					"data": map[string]interface{}{
						"getMembershipByLevel": patron,
					},
				}

				result, err := json.Marshal(response)
				if err != nil {
					log.Fatalf("Error marshalling patron to JSON: %v", err)
				}

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

			case "getMembershipByPatronId":
				patron_id, _ := data.Variables["patron_id"].(string)

				patron, ResolverErr := resolver.Query().GetMembershipByPatronID(ctx, patron_id)
				if ResolverErr != nil {
					log.Fatalf("Get Patron By ID Resolver err: %v", ResolverErr)
				}

				response := map[string]interface{}{
					"data": map[string]interface{}{
						"getMembershipByPatronId": patron,
					},
				}

				result, err := json.Marshal(response)
				if err != nil {
					log.Fatalf("Error marshalling patron to JSON: %v", err)
				}

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

			case "getViolationByPatronId":
				patron_id, _ := data.Variables["patron_id"].(string)

				patron, ResolverErr := resolver.Query().GetViolationByPatronID(ctx, patron_id)
				if ResolverErr != nil {
					log.Fatalf("Get Patron By ID Resolver err: %v", ResolverErr)
				}

				response := map[string]interface{}{
					"data": map[string]interface{}{
						"getViolationByPatronId": patron,
					},
				}

				result, err := json.Marshal(response)
				if err != nil {
					log.Fatalf("Error marshalling patron to JSON: %v", err)
				}

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

			case "getViolationByType":
				violation_type, _ := data.Variables["violation_type"].(string)

				patron, ResolverErr := resolver.Query().GetViolationByType(ctx, model.ViolationType(violation_type))
				if ResolverErr != nil {
					log.Fatalf("Get Patron By ID Resolver err: %v", ResolverErr)
				}

				response := map[string]interface{}{
					"data": map[string]interface{}{
						"getViolationByType": patron,
					},
				}

				result, err := json.Marshal(response)
				if err != nil {
					log.Fatalf("Error marshalling patron to JSON: %v", err)
				}

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

			case "getPatronStatusByType":
				status, _ := data.Variables["status"].(string)

				patron, ResolverErr := resolver.Query().GetPatronStatusByType(ctx, model.Status(status))
				if ResolverErr != nil {
					log.Fatalf("Get Patron By ID Resolver err: %v", ResolverErr)
				}

				response := map[string]interface{}{
					"data": map[string]interface{}{
						"getPatronStatusByType": patron,
					},
				}

				result, err := json.Marshal(response)
				if err != nil {
					log.Fatalf("Error marshalling patron to JSON: %v", err)
				}

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
				log.Printf("Unknown resolver: %s", data.RequestedResolver)
			}

		case <-sigChan:
			log.Println("Shutting down server...")
			return
		}
	}

}
