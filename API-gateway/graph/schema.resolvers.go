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
	"net/http"

	"github.com/Cat6utpcableclarke/API-gateway/graph/model"
)

const bookServiceURL = "http://localhost:8080/query" // Replace with your bookService URL

// AddBook is the resolver for the addBook field.
func (r *mutationResolver) AddBook(ctx context.Context, title string, authorName string, datePublished string, description string) (*model.Book, error) {
	query := `
        mutation AddBook($title: String!, $authorName: String!, $datePublished: String!, $description: String!) {
            addBook(title: $title, authorName: $authorName, datePublished: $datePublished, description: $description) {
                id
                title
                authorName
                datePublished
                description
            }
        }
    `

	variables := map[string]interface{}{
		"title":         title,
		"authorName":    authorName,
		"datePublished": datePublished,
		"description":   description,
	}

	resp, err := forwardRequest(ctx, query, variables)
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

	resp, err := forwardRequest(ctx, query, nil)
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

	resp, err := forwardRequest(ctx, query, variables)
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

	resp, err := forwardRequest(ctx, query, variables)
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

// BookAdded is the resolver for the bookAdded field.
func (r *subscriptionResolver) BookAdded(ctx context.Context) (<-chan *model.Book, error) {
	// Create a channel to forward subscription events
	bookChannel := make(chan *model.Book, 1)

	// Simulate forwarding the subscription to the bookService
	go func() {
		defer close(bookChannel)

		// Example: Simulate receiving a new book from the bookService
		newBook := &model.Book{
			ID:            "1",
			Title:         "New Book Title",
			AuthorName:    "Author Name",
			DatePublished: "2025-04-13",
			Description:   "A description of the new book.",
		}

		// Send the new book to the client
		bookChannel <- newBook
	}()

	return bookChannel, nil
}

// Helper function to forward requests to the bookService
func forwardRequest(ctx context.Context, query string, variables map[string]interface{}) ([]byte, error) {
	body := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", bookServiceURL, bytes.NewBuffer(jsonBody))
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

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
