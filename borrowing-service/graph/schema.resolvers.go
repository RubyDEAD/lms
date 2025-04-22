package graph

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/RubyDEAD/lms/borrowing-service/graph/model"
	"github.com/RubyDEAD/lms/borrowing-service/services"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// BorrowBook implements the borrowBook mutation
func (r *mutationResolver) BorrowBook(ctx context.Context, bookID string, patronID string) (*model.BorrowRecord, error) {
	// // Check book availability in Supabase
	// var book struct {
	// 	Available bool `json:"available"`
	// }
	// err := r.Supabase.DB.From("books").Select("available").Eq("id", bookID).Execute(&book)
	// if err != nil {
	// 	return nil, errors.New("failed to fetch a single book record")
	// }
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to check book availability: %v", err)
	// }
	// if !book.Available {
	// 	return nil, errors.New("book not available for borrowing")
	// }
	available, conn, bookCopyID, err := services.CheckAvailability(bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to check book availability: %v", err)
	}
	// Ensure the RabbitMQ connection is closed only if it is not already closed
	defer func() {
		if conn != nil && !conn.IsClosed() {
			conn.Close()
		}
	}()

	// If the book is not available, return an error
	if !available || bookCopyID == "0" {
		return nil, errors.New("no available book copies for borrowing")
	}
	err = services.SendUpdateRequest(conn, bookCopyID, "Borrowed")
	if err != nil {
		return nil, fmt.Errorf("failed to update book copy status: %v", err)
	}
	// Start transaction
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	// Create record
	now := time.Now()
	record := &model.BorrowRecord{
		ID:           uuid.New().String(),
		BookID:       bookID,
		PatronID:     patronID,
		BorrowedAt:   now.Format(time.RFC3339),
		DueDate:      now.AddDate(0, 0, 14).Format(time.RFC3339),
		RenewalCount: 0,
		Status:       model.BorrowStatusActive,
	}

	// Insert into PostgreSQL
	const insertQuery = `
		INSERT INTO borrow_records 
		(id, book_id, patron_id, borrowed_at, due_date, renewal_count, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err = tx.Exec(ctx, insertQuery,
		record.ID,
		record.BookID,
		record.PatronID,
		record.BorrowedAt,
		record.DueDate,
		record.RenewalCount,
		record.Status,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create borrow record: %v", err)
	}

	// Update book status in Supabase
	// err = r.Supabase.DB.From("books").
	// 	Update(map[string]interface{}{"available": false}).
	// 	Eq("id", bookID).
	// 	Execute(nil)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to update book status: %v", err)
	// }

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return record, nil
}

// ReturnBook implements the returnBook mutation
func (r *mutationResolver) ReturnBook(ctx context.Context, recordID string) (*model.BorrowRecord, error) {
	// Start transaction
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	// Get the borrow record
	var record model.BorrowRecord
	const getQuery = `
		SELECT id, book_id, patron_id, borrowed_at, due_date, returned_at, renewal_count, status 
		FROM borrow_records 
		WHERE id = $1
		FOR UPDATE
	`
	err = tx.QueryRow(ctx, getQuery, recordID).Scan(
		&record.ID,
		&record.BookID,
		&record.PatronID,
		&record.BorrowedAt,
		&record.DueDate,
		&record.ReturnedAt,
		&record.RenewalCount,
		&record.Status,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("borrow record not found")
		}
		return nil, fmt.Errorf("failed to get borrow record: %v", err)
	}

	if record.Status == model.BorrowStatusReturned {
		return nil, errors.New("book already returned")
	}

	// Update the record
	returnedAt := time.Now().Format(time.RFC3339)
	const updateQuery = `
		UPDATE borrow_records 
		SET returned_at = $1, status = $2
		WHERE id = $3
		RETURNING id, book_id, patron_id, borrowed_at, due_date, returned_at, renewal_count, status
	`
	err = tx.QueryRow(ctx, updateQuery, returnedAt, model.BorrowStatusReturned, recordID).Scan(
		&record.ID,
		&record.BookID,
		&record.PatronID,
		&record.BorrowedAt,
		&record.DueDate,
		&record.ReturnedAt,
		&record.RenewalCount,
		&record.Status,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update borrow record: %v", err)
	}

	// Update book availability in Supabase
	err = r.Supabase.DB.From("books").
		Update(map[string]interface{}{"available": true}).
		Eq("id", record.BookID).
		Execute(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to update book status: %v", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return &record, nil
}

// RenewLoan implements the renewLoan mutation
func (r *mutationResolver) RenewLoan(ctx context.Context, recordID string) (model.RenewLoanResult, error) {
	// Start transaction
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	// Get current record
	var record model.BorrowRecord
	const getQuery = `
		SELECT id, book_id, patron_id, borrowed_at, due_date, returned_at, renewal_count, status 
		FROM borrow_records 
		WHERE id = $1
		FOR UPDATE
	`
	err = tx.QueryRow(ctx, getQuery, recordID).Scan(
		&record.ID,
		&record.BookID,
		&record.PatronID,
		&record.BorrowedAt,
		&record.DueDate,
		&record.ReturnedAt,
		&record.RenewalCount,
		&record.Status,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("borrow record not found")
		}
		return nil, fmt.Errorf("failed to get borrow record: %v", err)
	}

	// Validate
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

	// Check for reservations in Supabase
	var resCount int
	err = r.Supabase.DB.From("reservations").
		Select("count").
		Eq("book_id", record.BookID).
		Eq("status", string(model.ReservationStatusPending)).
		Execute(&resCount)
	if err == nil && resCount > 0 {
		return &model.RenewalError{
			Code:    model.RenewalErrorCodeItemReserved,
			Message: "book is reserved by another patron",
		}, nil
	}

	// Update record
	newDueDate := time.Now().AddDate(0, 0, 14).Format(time.RFC3339)
	previousDueDate := record.DueDate
	const updateQuery = `
		UPDATE borrow_records 
		SET previous_due_date = $1, due_date = $2, renewal_count = $3, status = $4
		WHERE id = $5
		RETURNING id, book_id, patron_id, borrowed_at, due_date, returned_at, renewal_count, status
	`
	err = tx.QueryRow(ctx, updateQuery,
		previousDueDate,
		newDueDate,
		record.RenewalCount+1,
		model.BorrowStatusRenewed,
		recordID,
	).Scan(
		&record.ID,
		&record.BookID,
		&record.PatronID,
		&record.BorrowedAt,
		&record.DueDate,
		&record.ReturnedAt,
		&record.RenewalCount,
		&record.Status,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update borrow record: %v", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return &record, nil
}

// ReserveBook implements the reserveBook mutation
func (r *mutationResolver) ReserveBook(ctx context.Context, bookID string, patronID string) (*model.Reservation, error) {
	// Check book availability in Supabase
	var book struct {
		Available bool `json:"available"`
	}
	err := r.Supabase.DB.From("books").Select("available").Eq("id", bookID).Execute(&book)
	if err != nil {
		return nil, fmt.Errorf("failed to check book availability: %v", err)
	}
	if book.Available {
		return nil, errors.New("book is currently available, cannot reserve")
	}

	// Start transaction
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	// Create reservation
	now := time.Now()
	reservation := &model.Reservation{
		ID:         uuid.New().String(),
		BookID:     bookID,
		PatronID:   patronID,
		ReservedAt: now.Format(time.RFC3339),
		ExpiresAt:  now.AddDate(0, 0, 3).Format(time.RFC3339),
		Status:     model.ReservationStatusPending,
	}

	// Insert into PostgreSQL
	const insertQuery = `
		INSERT INTO reservations 
		(id, book_id, patron_id, reserved_at, expires_at, status)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err = tx.Exec(ctx, insertQuery,
		reservation.ID,
		reservation.BookID,
		reservation.PatronID,
		reservation.ReservedAt,
		reservation.ExpiresAt,
		reservation.Status,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create reservation: %v", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return reservation, nil
}

// CancelReservation implements the cancelReservation mutation
func (r *mutationResolver) CancelReservation(ctx context.Context, id string) (bool, error) {
	// Start transaction
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	// Get the reservation
	var reservation model.Reservation
	const getQuery = `
		SELECT id, book_id, patron_id, reserved_at, expires_at, status 
		FROM reservations 
		WHERE id = $1
		FOR UPDATE
	`
	err = tx.QueryRow(ctx, getQuery, id).Scan(
		&reservation.ID,
		&reservation.BookID,
		&reservation.PatronID,
		&reservation.ReservedAt,
		&reservation.ExpiresAt,
		&reservation.Status,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, errors.New("reservation not found")
		}
		return false, fmt.Errorf("failed to get reservation: %v", err)
	}

	if reservation.Status != model.ReservationStatusPending {
		return false, errors.New("only pending reservations can be cancelled")
	}

	// Update reservation
	const updateQuery = `
		UPDATE reservations 
		SET status = $1
		WHERE id = $2
	`
	_, err = tx.Exec(ctx, updateQuery, model.ReservationStatusCancelled, id)
	if err != nil {
		return false, fmt.Errorf("failed to cancel reservation: %v", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return false, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return true, nil
}

// FulfillReservation implements the fulfillReservation mutation
func (r *mutationResolver) FulfillReservation(ctx context.Context, id string) (*model.Reservation, error) {
	// Start transaction
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	// Get the reservation
	var reservation model.Reservation
	const getQuery = `
		SELECT id, book_id, patron_id, reserved_at, expires_at, status 
		FROM reservations 
		WHERE id = $1
		FOR UPDATE
	`
	err = tx.QueryRow(ctx, getQuery, id).Scan(
		&reservation.ID,
		&reservation.BookID,
		&reservation.PatronID,
		&reservation.ReservedAt,
		&reservation.ExpiresAt,
		&reservation.Status,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("reservation not found")
		}
		return nil, fmt.Errorf("failed to get reservation: %v", err)
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
		// Update reservation status to expired
		_, err = tx.Exec(ctx, `
			UPDATE reservations 
			SET status = $1
			WHERE id = $2
		`, model.ReservationStatusExpired, id)
		if err != nil {
			return nil, fmt.Errorf("failed to mark reservation as expired: %v", err)
		}
		return nil, errors.New("reservation has expired")
	}

	// Create borrow record
	_, err = r.BorrowBook(ctx, reservation.BookID, reservation.PatronID)
	if err != nil {
		return nil, fmt.Errorf("failed to create borrow record: %v", err)
	}

	// Update reservation status
	const updateQuery = `
		UPDATE reservations 
		SET status = $1
		WHERE id = $2
		RETURNING id, book_id, patron_id, reserved_at, expires_at, status
	`
	err = tx.QueryRow(ctx, updateQuery, model.ReservationStatusFulfilled, id).Scan(
		&reservation.ID,
		&reservation.BookID,
		&reservation.PatronID,
		&reservation.ReservedAt,
		&reservation.ExpiresAt,
		&reservation.Status,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update reservation: %v", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return &reservation, nil
}

// Query resolvers...

func (r *queryResolver) BorrowRecords(ctx context.Context, patronID *string, bookID *string, status *model.BorrowStatus) ([]*model.BorrowRecord, error) {
	query := `SELECT id, book_id, patron_id, borrowed_at, due_date, returned_at, renewal_count, status FROM borrow_records WHERE 1=1`
	args := []interface{}{}

	if patronID != nil {
		query += " AND patron_id = $1"
		args = append(args, *patronID)
	}
	if bookID != nil {
		if len(args) == 0 {
			query += " AND book_id = $1"
		} else {
			query += " AND book_id = $2"
		}
		args = append(args, *bookID)
	}
	if status != nil {
		pos := len(args) + 1
		query += fmt.Sprintf(" AND status = $%d", pos)
		args = append(args, string(*status))
	}

	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query borrow records: %v", err)
	}
	defer rows.Close()

	var records []*model.BorrowRecord
	for rows.Next() {
		var record model.BorrowRecord
		err := rows.Scan(
			&record.ID,
			&record.BookID,
			&record.PatronID,
			&record.BorrowedAt,
			&record.DueDate,
			&record.ReturnedAt,
			&record.RenewalCount,
			&record.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan borrow record: %v", err)
		}
		records = append(records, &record)
	}

	return records, nil
}

func (r *queryResolver) Reservations(ctx context.Context, patronID *string, bookID *string, status *model.ReservationStatus) ([]*model.Reservation, error) {
	query := `SELECT id, book_id, patron_id, reserved_at, expires_at, status FROM reservations WHERE 1=1`
	args := []interface{}{}

	if patronID != nil {
		query += " AND patron_id = $1"
		args = append(args, *patronID)
	}
	if bookID != nil {
		if len(args) == 0 {
			query += " AND book_id = $1"
		} else {
			query += " AND book_id = $2"
		}
		args = append(args, *bookID)
	}
	if status != nil {
		pos := len(args) + 1
		query += fmt.Sprintf(" AND status = $%d", pos)
		args = append(args, string(*status))
	}

	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query reservations: %v", err)
	}
	defer rows.Close()

	var reservations []*model.Reservation
	for rows.Next() {
		var res model.Reservation
		err := rows.Scan(
			&res.ID,
			&res.BookID,
			&res.PatronID,
			&res.ReservedAt,
			&res.ExpiresAt,
			&res.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan reservation: %v", err)
		}
		reservations = append(reservations, &res)
	}

	return reservations, nil
}

func (r *queryResolver) OverdueRecords(ctx context.Context) ([]*model.BorrowRecord, error) {
	now := time.Now().Format(time.RFC3339)
	query := `
		SELECT id, book_id, patron_id, borrowed_at, due_date, returned_at, renewal_count, status 
		FROM borrow_records 
		WHERE status = $1 AND due_date < $2
	`

	rows, err := r.DB.Query(ctx, query, model.BorrowStatusActive, now)
	if err != nil {
		return nil, fmt.Errorf("failed to query overdue records: %v", err)
	}
	defer rows.Close()

	var records []*model.BorrowRecord
	for rows.Next() {
		var record model.BorrowRecord
		err := rows.Scan(
			&record.ID,
			&record.BookID,
			&record.PatronID,
			&record.BorrowedAt,
			&record.DueDate,
			&record.ReturnedAt,
			&record.RenewalCount,
			&record.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan borrow record: %v", err)
		}
		records = append(records, &record)
	}

	return records, nil
}

func (r *queryResolver) PatronBorrowHistory(ctx context.Context, patronID string) ([]*model.BorrowRecord, error) {
	query := `
		SELECT id, book_id, patron_id, borrowed_at, due_date, returned_at, renewal_count, status 
		FROM borrow_records 
		WHERE patron_id = $1
		ORDER BY borrowed_at DESC
	`

	rows, err := r.DB.Query(ctx, query, patronID)
	if err != nil {
		return nil, fmt.Errorf("failed to query patron borrow history: %v", err)
	}
	defer rows.Close()

	var records []*model.BorrowRecord
	for rows.Next() {
		var record model.BorrowRecord
		err := rows.Scan(
			&record.ID,
			&record.BookID,
			&record.PatronID,
			&record.BorrowedAt,
			&record.DueDate,
			&record.ReturnedAt,
			&record.RenewalCount,
			&record.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan borrow record: %v", err)
		}
		records = append(records, &record)
	}

	return records, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
