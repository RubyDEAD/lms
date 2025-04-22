package graph

import (
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	supa "github.com/nedpals/supabase-go"
)

// Resolver handles GraphQL queries and mutations for borrowing service
type Resolver struct {
	DB       *pgx.Conn
	Supabase *supa.Client
}

// NewResolver creates a new resolver instance
func NewResolver(db *pgx.Conn, sb *supa.Client) *Resolver {
	return &Resolver{
		DB:       db,
		Supabase: sb,
	}
}

// BorrowRecord represents a book borrowing record
type BorrowRecord struct {
	ID           string             `json:"id"`
	BookID       string             `json:"book_id"`
	PatronID     string             `json:"patron_id"`
	BorrowedAt   time.Time          `json:"borrowed_at"`
	DueDate      time.Time          `json:"due_date"`
	ReturnedAt   pgtype.Timestamptz `json:"returned_at"`
	RenewalCount int                `json:"renewal_count"`
	Status       string             `json:"status"`
}

type BorrowRecordFilter struct {
	PatronID *string
	BookID   *string
	Status   *string
}
