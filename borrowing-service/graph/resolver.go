package graph

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// Resolver handles GraphQL queries and mutations
type Resolver struct {
	DB *pgx.Conn
}

// NewResolver creates a new resolver instance
func NewResolver(db *pgx.Conn) *Resolver {
	return &Resolver{DB: db}
}

// BeginTx starts a new database transaction
func (r *Resolver) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.DB.Begin(ctx)
}

// HealthCheck verifies database connection
func (r *Resolver) HealthCheck(ctx context.Context) error {
	return r.DB.Ping(ctx)
}

// Example query resolver
func (r *Resolver) GetBookByID(ctx context.Context, id string) (*Book, error) {
	const query = `SELECT id, title, author FROM books WHERE id = $1`

	var book Book
	err := r.DB.QueryRow(ctx, query, id).Scan(
		&book.ID,
		&book.Title,
		&book.Author,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("book not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &book, nil
}

// Example mutation resolver with transaction
func (r *Resolver) CreateBook(ctx context.Context, input NewBook) (*Book, error) {
	tx, err := r.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	const query = `
		INSERT INTO books (title, author) 
		VALUES ($1, $2)
		RETURNING id, title, author
	`

	var book Book
	err = tx.QueryRow(ctx, query, input.Title, input.Author).Scan(
		&book.ID,
		&book.Title,
		&book.Author,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create book: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &book, nil
}

// Book represents our data model
type Book struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

// NewBook represents input for creating a book
type NewBook struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}
