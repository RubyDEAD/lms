package graph

import (
	"fmt"
	"sync"
	"time"

	"github.com/RubyDEAD/lms/borrowing-service/graph/model"
	"github.com/jackc/pgx/v5"

	supa "github.com/nedpals/supabase-go"
)

// Resolver handles GraphQL queries and mutations for borrowing service
type Resolver struct {
	DB                            *pgx.Conn
	Supabase                      *supa.Client
	mutex                         sync.Mutex
	reservationCreatedChannels    map[string][]chan *model.Reservation
	reservedBookAvailableChannels map[string][]chan *model.Reservation
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
	ID              string     `json:"id"`
	BookID          string     `json:"book_id"`
	PatronID        string     `json:"patron_id"`
	BorrowedAt      time.Time  `json:"borrowed_at"`
	DueDate         time.Time  `json:"due_date"`
	ReturnedAt      *time.Time `json:"returned_at"`
	PreviousDueDate *time.Time `json:"previous_due_date"`
	RenewalCount    int        `json:"renewal_count"`
	Status          string     `json:"status"`
	BookCopyID      int        `json:"book_copy_id"`
}

type BorrowRecordFilter struct {
	PatronID *string
	BookID   *string
	Status   *string
}

type FulfillmentError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *FulfillmentError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func convertToModel(record BorrowRecord) model.BorrowRecord {
	borrowedAtStr := record.BorrowedAt.Format(time.RFC3339)
	dueDateStr := record.DueDate.Format(time.RFC3339)

	var returnedAtStr *string
	if record.ReturnedAt != nil {
		str := record.ReturnedAt.Format(time.RFC3339)
		returnedAtStr = &str
	}

	var prevDueDateStr *string
	if record.PreviousDueDate != nil {
		str := record.PreviousDueDate.Format(time.RFC3339)
		prevDueDateStr = &str
	}

	return model.BorrowRecord{
		ID:              record.ID,
		BookID:          record.BookID,
		PatronID:        record.PatronID,
		BorrowedAt:      borrowedAtStr,
		DueDate:         dueDateStr,
		ReturnedAt:      returnedAtStr,
		RenewalCount:    int32(record.RenewalCount),
		Status:          model.BorrowStatus(record.Status),
		BookCopyID:      int32(record.BookCopyID),
		PreviousDueDate: prevDueDateStr,
	}
}
