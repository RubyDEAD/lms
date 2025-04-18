package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.70

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/Cat6utpcableclarke/API-gateway/graph/model"
	"github.com/coder/websocket"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

const bookServiceURL = "http://localhost:8080/query"
const patronServiceQueue = "patron-service-queue"

// HELPER FUNCTION - FORWARDS REQUESTS FROM THE GATEWAY TO THE INDIVIDUAL SERVICES (VERY IMPORTANT)
func forwardRequest(ctx context.Context, query string, variables map[string]interface{}, serviceURL string) ([]byte, error) {
	body := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", serviceURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

// HELPER FUNCTION
func SubscribeToBookAdded(ctx context.Context, out chan<- *model.Book) error {
	c, _, err := websocket.Dial(ctx, "ws://localhost:8080/query", nil)
	if err != nil {
		return err
	}

	// connection_init
	if err := c.Write(ctx, websocket.MessageText, []byte(`{"type":"connection_init"}`)); err != nil {
		return err
	}
	_, _, _ = c.Read(ctx) // wait for ack

	// send subscription
	payload := map[string]interface{}{
		"id":   "1",
		"type": "start",
		"payload": map[string]interface{}{
			"query": `subscription { bookAdded { id title author_name date_published description } }`,
		},
	}
	msg, _ := json.Marshal(payload)
	if err := c.Write(ctx, websocket.MessageText, msg); err != nil {
		return err
	}

	go func() {
		for {
			_, data, err := c.Read(ctx)
			if err != nil {
				log.Println("WebSocket read failed:", err)
				close(out)
				return
			}

			var msg struct {
				Type    string `json:"type"`
				Payload struct {
					Data struct {
						BookAdded *model.Book `json:"bookAdded"`
					} `json:"data"`
				} `json:"payload"`
			}
			if err := json.Unmarshal(data, &msg); err != nil {
				log.Println("Failed to unmarshal message:", err)
				continue
			}

			if msg.Type == "data" {
				out <- msg.Payload.Data.BookAdded
			}
		}
	}()

	return nil
}

func forwardRequestMQ(queue string, variables map[string]interface{}, requestedResolver string) ([]byte, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}
	defer ch.Close()

	// Declare the request and reply queues
	_, err = ch.QueueDeclare(
		queue,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue (%s): %w", queue, err)
	}

	replyQueue, err := ch.QueueDeclare(
		"api-gateway-queue",
		false,
		true,
		true,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare reply queue: %w", err)
	}

	corrID := uuid.New().String()

	body := map[string]interface{}{
		"variables":         variables,
		"requestedResolver": requestedResolver,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Publish request
	err = ch.Publish(
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: corrID,
			ReplyTo:       replyQueue.Name,
			Body:          jsonBody,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to publish message: %w", err)
	}

	// Consume from reply queue
	msgs, err := ch.Consume(
		replyQueue.Name,
		"",
		true,  // auto-ack
		true,  // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to consume reply: %w", err)
	}

	// Wait for a response with the matching correlation ID
	timeout := time.After(5 * time.Second)
	for {
		select {
		case msg := <-msgs:
			if msg.CorrelationId == corrID {
				log.Printf("body: %s", msg.Body)
				log.Printf("id: %s", msg.CorrelationId)
				return msg.Body, nil
			}
		case <-timeout:
			return nil, fmt.Errorf("timeout waiting for reply from service")
		}
	}
}

// AddBook is the resolver for the addBook field.
func (r *mutationResolver) AddBook(ctx context.Context, title string, authorName string, datePublished string, description string) (*model.Book, error) {
	query := `
        mutation AddBook($title: String!, $author_name: String!, $datePublished: String!, $description: String!) {
            addBook(title: $title, author_name: $author_name, datePublished: $datePublished, description: $description) {
                id
                title
                author_name
                date_published
                description
            }
        }
    `

	variables := map[string]interface{}{
		"title":         title,
		"author_name":   authorName,
		"datePublished": datePublished,
		"description":   description,
	}

	resp, err := forwardRequest(ctx, query, variables, bookServiceURL)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data struct {
			AddBook *model.Book `json:"addBook"`
		} `json:"data"`
	}

	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return result.Data.AddBook, nil
}

// CreatePatron is the resolver for the createPatron field.
func (r *mutationResolver) CreatePatron(ctx context.Context, firstName string, lastName string, phoneNumber string) (*model.Patron, error) {
	variables := map[string]interface{}{
		"firstName":   firstName,
		"lastName":    lastName,
		"phoneNumber": phoneNumber,
	}

	//resp, err := forwardRequest(ctx, query, variables, patronServiceURL)
	resp, err := forwardRequestMQ(patronServiceQueue, variables, "createPatron")
	if err != nil {
		return nil, fmt.Errorf("failed to forward request: %v", err)
	}

	var result struct {
		Data struct {
			CreatePatron *model.Patron `json:"createPatron"`
		} `json:"data"`
	}
	// log.Printf("in createpatron: %s", resp)
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	// log.Printf("check unmarshal: %+v", result.Data.CreatePatron)

	return result.Data.CreatePatron, nil
}

// UpdatePatron is the resolver for the updatePatron field.
func (r *mutationResolver) UpdatePatron(ctx context.Context, patronID string, firstName *string, lastName *string, phoneNumber *string) (*model.Patron, error) {
	variables := map[string]interface{}{
		"patron_id":   patronID,
		"firstName":   firstName,
		"lastName":    lastName,
		"phoneNumber": phoneNumber,
	}

	resp, err := forwardRequestMQ(patronServiceQueue, variables, "updatePatron")
	if err != nil {
		return nil, fmt.Errorf("failed to forward request: %v", err)
	}

	var result struct {
		Data struct {
			UpdatePatron *model.Patron `json:"updatePatron"`
		} `json:"data"`
	}

	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return result.Data.UpdatePatron, nil
}

// DeletePatronByID is the resolver for the deletePatronById field.
func (r *mutationResolver) DeletePatronByID(ctx context.Context, patronID string) (*model.Patron, error) {
	variables := map[string]interface{}{
		"patron_id": patronID,
	}

	resp, err := forwardRequestMQ(patronServiceQueue, variables, "deletePatronById")
	if err != nil {
		return nil, fmt.Errorf("failed to forward request: %v", err)
	}

	var result struct {
		Data struct {
			DeletePatronByID *model.Patron `json:"deletePatronById"`
		} `json:"data"`
	}

	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return result.Data.DeletePatronByID, nil
}

// UpdateMembershipByPatronID is the resolver for the updateMembershipByPatronId field.
func (r *mutationResolver) UpdateMembershipByPatronID(ctx context.Context, patronID string, level model.MembershipLevel) (*model.Membership, error) {
	variables := map[string]interface{}{
		"patron_id": patronID,
		"level":     level,
	}

	resp, err := forwardRequestMQ(patronServiceQueue, variables, "updateMembershipByPatronId")
	if err != nil {
		return nil, fmt.Errorf("failed to forward request: %v", err)
	}

	var result struct {
		Data struct {
			UpdateMembershipByPatronID *model.Membership `json:"updateMembershipByPatronId"`
		} `json:"data"`
	}

	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshall response: %v", err)
	}

	return result.Data.UpdateMembershipByPatronID, nil
}

// UpdateMembershipByMembershipID is the resolver for the updateMembershipByMembershipId field.
func (r *mutationResolver) UpdateMembershipByMembershipID(ctx context.Context, membershipID string, level model.MembershipLevel) (*model.Membership, error) {
	variables := map[string]interface{}{
		"membership_id": membershipID,
		"level":         level,
	}

	resp, err := forwardRequestMQ(patronServiceQueue, variables, "updateMembershipByMembershipId")
	if err != nil {
		return nil, fmt.Errorf("failed to forward request: %v", err)
	}

	var result struct {
		Data struct {
			UpdateMembershipByMembershipID *model.Membership `json:"updateMembershipByMembershipId"`
		} `json:"data"`
	}

	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshall response: %v", err)
	}

	return result.Data.UpdateMembershipByMembershipID, nil
}

// UpdatePatronStatus is the resolver for the updatePatronStatus field.
func (r *mutationResolver) UpdatePatronStatus(ctx context.Context, patronID string, warningCount *int32, unpaidFees *float64, patronStatus *model.Status) (*model.PatronStatus, error) {
	variables := map[string]interface{}{
		"patron_id":     patronID,
		"warning_count": warningCount,
		"unpaid_fees":   unpaidFees,
		"status":        patronStatus,
	}

	resp, err := forwardRequestMQ(patronServiceQueue, variables, "updatePatronStatus")
	if err != nil {
		return nil, fmt.Errorf("failed to forward request: %v", err)
	}

	var result struct {
		Data struct {
			UpdatePatronStatus *model.PatronStatus `json:"updatePatronStatus"`
		} `json:"data"`
	}

	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshall response: %v", err)
	}

	return result.Data.UpdatePatronStatus, nil
}

// AddViolation is the resolver for the addViolation field.
func (r *mutationResolver) AddViolation(ctx context.Context, patronID string, violationType model.ViolationType, violationInfo string) (*model.ViolationRecord, error) {
	variables := map[string]interface{}{
		"patron_id":      patronID,
		"violation_type": violationType,
		"violation_info": violationInfo,
	}

	resp, err := forwardRequestMQ(patronServiceQueue, variables, "addViolation")
	if err != nil {
		return nil, fmt.Errorf("failed to forward request: %v", err)
	}

	var result struct {
		Data struct {
			AddViolation *model.ViolationRecord `json:"addViolation"`
		} `json:"data"`
	}

	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshall response: %v", err)
	}

	return result.Data.AddViolation, nil
}

// UpdateViolationStatus is the resolver for the updateViolationStatus field.
func (r *mutationResolver) UpdateViolationStatus(ctx context.Context, violationID string, violationStatus model.ViolationStatus) (*model.ViolationRecord, error) {
	variables := map[string]interface{}{
		"violation_id":   violationID,
		"violation_type": violationStatus,
	}

	resp, err := forwardRequestMQ(patronServiceQueue, variables, "updateViolationStatus")
	if err != nil {
		return nil, fmt.Errorf("failed to forward request: %v", err)
	}

	var result struct {
		Data struct {
			UpdateViolationStatus *model.ViolationRecord `json:"updateViolationStatus"`
		} `json:"data"`
	}

	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshall response: %v", err)
	}

	return result.Data.UpdateViolationStatus, nil
}

// GetBooks is the resolver for the getBooks field.
func (r *queryResolver) GetBooks(ctx context.Context) ([]*model.Book, error) {
	query := `
        query {
            getBooks {
                id
                title
                author_name
                date_published
                description
            }
        }
    `

	resp, err := forwardRequest(ctx, query, nil, bookServiceURL)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data struct {
			GetBooks []*model.Book `json:"getBooks"`
		} `json:"data"`
	}

	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return result.Data.GetBooks, nil
}

// GetBookByID is the resolver for the getBookById field.
func (r *queryResolver) GetBookByID(ctx context.Context, id string) (*model.Book, error) {
	query := `
        query GetBookById($id: String!) {
            getBookById(id: $id) {
                id
                title
                author_name
                date_published
                description
            }
        }
    `

	variables := map[string]interface{}{
		"id": id,
	}

	resp, err := forwardRequest(ctx, query, variables, bookServiceURL)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data struct {
			GetBookById *model.Book `json:"getBookById"`
		} `json:"data"`
	}

	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return result.Data.GetBookById, nil
}

// GetBookCopiesByID is the resolver for the getBookCopiesById field.
func (r *queryResolver) GetBookCopiesByID(ctx context.Context, id string) ([]*model.BookCopies, error) {
	query := `
        query GetBookCopiesById($id: String!) {
            getBookCopiesById(id: $id) {
                id
                book_id
                book_status
            }
        }
    `

	variables := map[string]interface{}{
		"id": id,
	}

	resp, err := forwardRequest(ctx, query, variables, bookServiceURL)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data struct {
			GetBookCopiesById []*model.BookCopies `json:"getBookCopiesById"`
		} `json:"data"`
	}

	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return result.Data.GetBookCopiesById, nil
}

// GetPatronByID is the resolver for the getPatronById field.
func (r *queryResolver) GetPatronByID(ctx context.Context, patronID string) (*model.Patron, error) {
	variables := map[string]interface{}{
		"patron_id": patronID,
	}

	resp, err := forwardRequestMQ(patronServiceQueue, variables, "getPatronById")
	if err != nil {
		return nil, fmt.Errorf("failed to forward request: %v", err)
	}

	var result struct {
		Data struct {
			GetPatronByID *model.Patron `json:"getPatronById"`
		} `json:"data"`
	}

	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshall response: %v", err)
	}

	return result.Data.GetPatronByID, nil
}

// GetAllPatrons is the resolver for the getAllPatrons field.
func (r *queryResolver) GetAllPatrons(ctx context.Context) ([]*model.Patron, error) {
	variables := map[string]interface{}{}

	resp, err := forwardRequestMQ(patronServiceQueue, variables, "getAllPatrons")
	if err != nil {
		return nil, fmt.Errorf("failed to forward request: %v", err)
	}

	var result struct {
		Data struct {
			GetAllPatrons []*model.Patron `json:"getAllPatrons"`
		} `json:"data"`
	}

	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshall response: %v", err)
	}

	return result.Data.GetAllPatrons, nil
}

// GetMembershipByLevel is the resolver for the getMembershipByLevel field.
func (r *queryResolver) GetMembershipByLevel(ctx context.Context, level model.MembershipLevel) ([]*model.Membership, error) {
	variables := map[string]interface{}{
		"level": level,
	}

	resp, err := forwardRequestMQ(patronServiceQueue, variables, "getMembershipByLevel")
	if err != nil {
		return nil, fmt.Errorf("failed to forward request: %v", err)
	}

	var result struct {
		Data struct {
			GetMembershipByLevel []*model.Membership `json:"getMembershipByLevel"`
		} `json:"data"`
	}

	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshall response: %v", err)
	}

	return result.Data.GetMembershipByLevel, nil
}

// GetMembershipByPatronID is the resolver for the getMembershipByPatronId field.
func (r *queryResolver) GetMembershipByPatronID(ctx context.Context, patronID string) (*model.Membership, error) {
	variables := map[string]interface{}{
		"patron_id": patronID,
	}

	resp, err := forwardRequestMQ(patronServiceQueue, variables, "getMembershipByPatronId")
	if err != nil {
		return nil, fmt.Errorf("failed to forward request: %v", err)
	}

	var result struct {
		Data struct {
			GetMembershipByPatronID *model.Membership `json:"getMembershipByPatronId"`
		} `json:"data"`
	}

	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshall response: %v", err)
	}

	return result.Data.GetMembershipByPatronID, nil
}

// GetViolationByPatronID is the resolver for the getViolationByPatronId field.
func (r *queryResolver) GetViolationByPatronID(ctx context.Context, patronID string) ([]*model.ViolationRecord, error) {
	variables := map[string]interface{}{
		"patron_id": patronID,
	}

	resp, err := forwardRequestMQ(patronServiceQueue, variables, "getViolationByPatronId")
	if err != nil {
		return nil, fmt.Errorf("failed to forward request: %v", err)
	}

	var result struct {
		Data struct {
			GetViolationByPatronID []*model.ViolationRecord `json:"getViolationByPatronId"`
		} `json:"data"`
	}

	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshall response: %v", err)
	}

	return result.Data.GetViolationByPatronID, nil
}

// GetViolationByType is the resolver for the getViolationByType field.
func (r *queryResolver) GetViolationByType(ctx context.Context, violationType model.ViolationType) ([]*model.ViolationRecord, error) {
	variables := map[string]interface{}{
		"violation_type": violationType,
	}

	resp, err := forwardRequestMQ(patronServiceQueue, variables, "getViolationByType")
	if err != nil {
		return nil, fmt.Errorf("failed to forward request: %v", err)
	}

	var result struct {
		Data struct {
			GetViolationByType []*model.ViolationRecord `json:"getViolationByType"`
		} `json:"data"`
	}

	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshall response: %v", err)
	}

	return result.Data.GetViolationByType, nil
}

// GetPatronStatusByType is the resolver for the getPatronStatusByType field.
func (r *queryResolver) GetPatronStatusByType(ctx context.Context, patronStatus model.Status) ([]*model.PatronStatus, error) {
	variables := map[string]interface{}{
		"status": patronStatus,
	}

	resp, err := forwardRequestMQ(patronServiceQueue, variables, "getPatronStatusByType")
	if err != nil {
		return nil, fmt.Errorf("failed to forward request: %v", err)
	}

	var result struct {
		Data struct {
			GetPatronStatusByType []*model.PatronStatus `json:"getPatronStatusByType"`
		} `json:"data"`
	}

	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshall response: %v", err)
	}

	return result.Data.GetPatronStatusByType, nil
}

// BookAdded is the resolver for the bookAdded field.
func (r *subscriptionResolver) BookAdded(ctx context.Context) (<-chan *model.Book, error) {
	bookChan := make(chan *model.Book)
	err := SubscribeToBookAdded(ctx, bookChan)
	if err != nil {
		return nil, err
	}
	return bookChan, nil
}

// PatronCreated is the resolver for the patronCreated field.
func (r *subscriptionResolver) PatronCreated(ctx context.Context) (<-chan *model.Patron, error) {
	panic(fmt.Errorf("not implemented: PatronCreated - patronCreated"))
}

// OngoingViolations is the resolver for the ongoingViolations field.
func (r *subscriptionResolver) OngoingViolations(ctx context.Context) (<-chan *model.ViolationRecord, error) {
	panic(fmt.Errorf("not implemented: OngoingViolations - ongoingViolations"))
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
/*


 */
