package graph

import (
	"borrowing-service/graph/model"
	"context"
	"errors"
	"fmt"
	"time"
)

// Database interface (mock this for testing)
type DB interface {
	GetBookAvailability(ctx context.Context, bookID string) (bool, error)
	GetPatronStatus(ctx context.Context, patronID string) (string, error)
	CreateBorrowRecord(ctx context.Context, input model.BorrowRecord) (*model.BorrowRecord, error)
	UpdateBorrowRecord(ctx context.Context, recordID string, updates map[string]interface{}) (*model.BorrowRecord, error)
	GetBorrowRecord(ctx context.Context, recordID string) (*model.BorrowRecord, error)
	CreateReservation(ctx context.Context, input model.Reservation) (*model.Reservation, error)
	CancelReservation(ctx context.Context, id string) error
	// Add other required DB methods...
}

type Resolver struct {
	DB DB
	
	// Add other dependencies like bookServiceClient, patronServiceClient
}

// BorrowBook implements the borrow operation
func (r *mutationResolver) BorrowBook(ctx context.Context, bookID string, patronID string) (*model.BorrowRecord, error) {
	// 1. Check book availability
	available, err := r.DB.GetBookAvailability(ctx, bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to check book availability: %w", err)
	}
	if !available {
		return nil, errors.New("book is not available for borrowing")
	}

	// 2. Check patron status
	status, err := r.DB.GetPatronStatus(ctx, patronID)
	if err != nil {
		return nil, fmt.Errorf("failed to check patron status: %w", err)
	}
	if status != "ACTIVE" {
		return nil, fmt.Errorf("patron status is %s - cannot borrow", status)
	}

	// 3. Create borrow record
	now := time.Now().Format(time.RFC3339)
	dueDate := time.Now().AddDate(0, 0, 21).Format(time.RFC3339) // 3 weeks loan period

	record := model.BorrowRecord{
		BookID:      bookID,
		PatronID:    patronID,
		BorrowedAt:  now,
		DueDate:     dueDate,
		Status:      model.BorrowStatusActive,
		RenewalCount: 0,
	}

	return r.DB.CreateBorrowRecord(ctx, record)
}

// ReturnBook implements the return operation
func (r *mutationResolver) ReturnBook(ctx context.Context, recordID string) (*model.BorrowRecord, error) {
	// 1. Get the current record
	record, err := r.DB.GetBorrowRecord(ctx, recordID)
	if err != nil {
		return nil, fmt.Errorf("failed to find borrow record: %w", err)
	}

	// 2. Check if already returned
	if record.ReturnedAt != "" {
		return nil, errors.New("book already returned")
	}

	// 3. Update record
	updates := map[string]interface{}{
		"returnedAt": time.Now().Format(time.RFC3339),
		"status":     model.BorrowStatusReturned,
	}

	return r.DB.UpdateBorrowRecord(ctx, recordID, updates)
}

// RenewLoan implements the renewal operation
func (r *mutationResolver) RenewLoan(ctx context.Context, recordID string) (model.RenewLoanResult, error) {
	// 1. Get current record
	record, err := r.DB.GetBorrowRecord(ctx, recordID)
	if err != nil {
		return nil, fmt.Errorf("failed to find borrow record: %w", err)
	}

	// 2. Validate renewal conditions
	if record.ReturnedAt != "" {
		return model.RenewalError{
			Code:    model.RenewalErrorCodeLoanAlreadyReturned,
			Message: "Cannot renew returned items",
		}, nil
	}

	if record.RenewalCount >= 2 { // Max 2 renewals
		return model.RenewalError{
			Code:    model.RenewalErrorCodeMaxRenewalsReached,
			Message: "Maximum renewals reached",
		}, nil
	}

	// 3. Check if book has reservations
	// (Implement this check based on your DB)
	hasReservations, err := r.DB.BookHasReservations(ctx, record.BookID)
	if err != nil {
		return nil, fmt.Errorf("failed to check reservations: %w", err)
	}
	if hasReservations {
		return model.RenewalError{
			Code:    model.RenewalErrorCodeItemReserved,
			Message: "Cannot renew - item has reservations",
		}, nil
	}

	// 4. Process renewal
	newDueDate := time.Now().AddDate(0, 0, 21).Format(time.RFC3339) // Another 3 weeks
	updates := map[string]interface{}{
		"dueDate":        newDueDate,
		"previousDueDate": record.DueDate,
		"renewalCount":    record.RenewalCount + 1,
		"status":         model.BorrowStatusRenewed,
	}

	updated, err := r.DB.UpdateBorrowRecord(ctx, recordID, updates)
	if err != nil {
		return nil, fmt.Errorf("failed to renew loan: %w", err)
	}

	return updated, nil
}

// ReserveBook implements the reservation operation
func (r *mutationResolver) ReserveBook(ctx context.Context, bookID string, patronID string) (*model.Reservation, error) {
	// 1. Check patron eligibility
	status, err := r.DB.GetPatronStatus(ctx, patronID)
	if err != nil {
		return nil, fmt.Errorf("failed to check patron status: %w", err)
	}
	if status != "ACTIVE" {
		return nil, fmt.Errorf("patron status is %s - cannot reserve", status)
	}

	// 2. Create reservation
	now := time.Now().Format(time.RFC3339)
	expiresAt := time.Now().AddDate(0, 0, 7).Format(time.RFC3339) // 1 week reservation

	reservation := model.Reservation{
		BookID:     bookID,
		PatronID:   patronID,
		ReservedAt: now,
		ExpiresAt:  expiresAt,
		Status:     model.ReservationStatusPending,
	}

	return r.DB.CreateReservation(ctx, reservation)
}

// Query resolvers...

func (r *queryResolver) BorrowRecords(ctx context.Context, patronID *string, bookID *string, status *model.BorrowStatus) ([]*model.BorrowRecord, error) {
	filters := make(map[string]interface{})
	if patronID != nil {
		filters["patron_id"] = *patronID
	}
	if bookID != nil {
		filters["book_id"] = *bookID
	}
	if status != nil {
		filters["status"] = *status
	}
	
	// Implement actual DB query
	return r.DB.GetBorrowRecords(ctx, filters)
}

// ... implement other query resolvers similarly

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }