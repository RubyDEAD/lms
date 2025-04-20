package graph

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/RubyDEAD/lms/borrowing-service/graph/model"
)

// In-memory storage (for demonstration, replace with real DB in production)
var (
	borrowRecords = make(map[string]*model.BorrowRecord)
	reservations  = make(map[string]*model.Reservation)
	bookStatus    = make(map[string]bool) // true = available
	patronStatus  = make(map[string]bool) // true = can borrow
)

// BorrowBook implements the borrowBook mutation
func (r *mutationResolver) BorrowBook(ctx context.Context, bookID string, patronID string) (*model.BorrowRecord, error) {
	// Check book availability
	if available, exists := bookStatus[bookID]; !exists || !available {
		return nil, errors.New("book not available for borrowing")
	}

	// Check patron status
	if canBorrow, exists := patronStatus[patronID]; !exists || !canBorrow {
		return nil, errors.New("patron cannot borrow at this time")
	}

	// Create record
	now := time.Now()
	record := &model.BorrowRecord{
		ID:           fmt.Sprintf("borrow-%d", len(borrowRecords)+1),
		BookID:       bookID,
		PatronID:     patronID,
		BorrowedAt:   now.Format(time.RFC3339),
		DueDate:      now.AddDate(0, 0, 14).Format(time.RFC3339), // 2 weeks
		RenewalCount: 0,
		Status:       model.BorrowStatusActive,
	}

	// Store record
	borrowRecords[record.ID] = record
	bookStatus[bookID] = false // Mark as borrowed

	return record, nil
}

// ReturnBook implements the returnBook mutation
func (r *mutationResolver) ReturnBook(ctx context.Context, recordID string) (*model.BorrowRecord, error) {
	record, exists := borrowRecords[recordID]
	if !exists {
		return nil, errors.New("borrow record not found")
	}

	if record.Status == model.BorrowStatusReturned {
		return nil, errors.New("book already returned")
	}

	// Update record
	formattedTime := time.Now().Format(time.RFC3339)
	record.ReturnedAt = &formattedTime
	record.Status = model.BorrowStatusReturned
	bookStatus[record.BookID] = true // Mark as available

	return record, nil
}

// RenewLoan implements the renewLoan mutation
func (r *mutationResolver) RenewLoan(ctx context.Context, recordID string) (model.RenewLoanResult, error) {
	record, exists := borrowRecords[recordID]
	if !exists {
		return nil, errors.New("borrow record not found")
	}

	// Validation checks
	if record.Status == model.BorrowStatusReturned {
		return &model.RenewalError{
			Code:    model.RenewalErrorCodeLoanAlreadyReturned,
			Message: "cannot renew already returned loan",
		}, nil
	}

	if record.RenewalCount >= 2 {
		return &model.RenewalError{
			Code:    model.RenewalErrorCodeMaxRenewalsReached,
			Message: "maximum renewals reached",
		}, nil
	}

	// Check for reservations
	for _, res := range reservations {
		if res.BookID == record.BookID && res.Status == model.ReservationStatusPending {
			return &model.RenewalError{
				Code:    model.RenewalErrorCodeItemReserved,
				Message: "book is reserved by another patron",
			}, nil
		}
	}

	// Update record
	record.PreviousDueDate = &record.DueDate
	record.DueDate = time.Now().AddDate(0, 0, 14).Format(time.RFC3339)
	record.RenewalCount++
	record.Status = model.BorrowStatusRenewed

	return record, nil
}

// ReserveBook implements the reserveBook mutation
func (r *mutationResolver) ReserveBook(ctx context.Context, bookID string, patronID string) (*model.Reservation, error) {
	// Check book availability
	if available, exists := bookStatus[bookID]; !exists || available {
		return nil, errors.New("book is currently available, cannot reserve")
	}

	// Check patron status
	if canBorrow, exists := patronStatus[patronID]; !exists || !canBorrow {
		return nil, errors.New("patron cannot make reservations")
	}

	// Create reservation
	now := time.Now()
	reservation := &model.Reservation{
		ID:         fmt.Sprintf("reserve-%d", len(reservations)+1),
		BookID:     bookID,
		PatronID:   patronID,
		ReservedAt: now.Format(time.RFC3339),
		ExpiresAt:  now.AddDate(0, 0, 3).Format(time.RFC3339), // 3 days
		Status:     model.ReservationStatusPending,
	}

	// Store reservation
	reservations[reservation.ID] = reservation

	return reservation, nil
}

// CancelReservation implements the cancelReservation mutation
func (r *mutationResolver) CancelReservation(ctx context.Context, id string) (bool, error) {
	reservation, exists := reservations[id]
	if !exists {
		return false, errors.New("reservation not found")
	}

	if reservation.Status != model.ReservationStatusPending {
		return false, errors.New("only pending reservations can be cancelled")
	}

	reservation.Status = model.ReservationStatusCancelled
	return true, nil
}

// FulfillReservation implements the fulfillReservation mutation
func (r *mutationResolver) FulfillReservation(ctx context.Context, id string) (*model.Reservation, error) {
	reservation, exists := reservations[id]
	if !exists {
		return nil, errors.New("reservation not found")
	}

	if reservation.Status != model.ReservationStatusPending {
		return nil, errors.New("only pending reservations can be fulfilled")
	}

	// Check expiration
	expiresAt, err := time.Parse(time.RFC3339, reservation.ExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("invalid expiration date: %v", err)
	}
	if time.Now().After(expiresAt) {
		reservation.Status = model.ReservationStatusExpired
		return nil, errors.New("reservation has expired")
	}

	// Create borrow record
	_, err = r.BorrowBook(ctx, reservation.BookID, reservation.PatronID)
	if err != nil {
		return nil, fmt.Errorf("failed to create borrow record: %v", err)
	}

	reservation.Status = model.ReservationStatusFulfilled
	return reservation, nil
}

// Query resolvers...

func (r *queryResolver) BorrowRecords(ctx context.Context, patronID *string, bookID *string, status *model.BorrowStatus) ([]*model.BorrowRecord, error) {
	var results []*model.BorrowRecord
	for _, record := range borrowRecords {
		match := true
		if patronID != nil && *patronID != record.PatronID {
			match = false
		}
		if bookID != nil && *bookID != record.BookID {
			match = false
		}
		if status != nil && *status != record.Status {
			match = false
		}
		if match {
			results = append(results, record)
		}
	}
	return results, nil
}

func (r *queryResolver) Reservations(ctx context.Context, patronID *string, bookID *string, status *model.ReservationStatus) ([]*model.Reservation, error) {
	var results []*model.Reservation
	for _, res := range reservations {
		match := true
		if patronID != nil && *patronID != res.PatronID {
			match = false
		}
		if bookID != nil && *bookID != res.BookID {
			match = false
		}
		if status != nil && *status != res.Status {
			match = false
		}
		if match {
			results = append(results, res)
		}
	}
	return results, nil
}

func (r *queryResolver) OverdueRecords(ctx context.Context) ([]*model.BorrowRecord, error) {
	var results []*model.BorrowRecord
	now := time.Now()
	for _, record := range borrowRecords {
		if record.Status == model.BorrowStatusActive {
			dueDate, err := time.Parse(time.RFC3339, record.DueDate)
			if err == nil && now.After(dueDate) {
				results = append(results, record)
			}
		}
	}
	return results, nil
}

func (r *queryResolver) PatronBorrowHistory(ctx context.Context, patronID string) ([]*model.BorrowRecord, error) {
	var results []*model.BorrowRecord
	for _, record := range borrowRecords {
		if record.PatronID == patronID {
			results = append(results, record)
		}
	}
	return results, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
/*
	func (r *mutationResolver) CreateTodo(ctx context.Context, input model.NewTodo) (*model.Todo, error) {
	panic(fmt.Errorf("not implemented: CreateTodo - createTodo"))
}
func (r *queryResolver) Todos(ctx context.Context) ([]*model.Todo, error) {
	panic(fmt.Errorf("not implemented: Todos - todos"))
}
*/
