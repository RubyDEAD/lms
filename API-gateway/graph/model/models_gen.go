// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
)

type RenewLoanResult interface {
	IsRenewLoanResult()
}

type Author struct {
	ID         int32  `json:"id"`
	AuthorName string `json:"author_name"`
}

type Book struct {
	ID            string  `json:"id"`
	Title         string  `json:"title"`
	AuthorName    string  `json:"author_name"`
	DatePublished string  `json:"date_published"`
	Description   string  `json:"description"`
	Image         *string `json:"image,omitempty"`
}

type BookCopies struct {
	ID            string `json:"id"`
	BookID        string `json:"book_id"`
	Title         string `json:"title"`
	AuthorName    string `json:"author_name"`
	DatePublished string `json:"date_published"`
	Description   string `json:"description"`
	BookStatus    string `json:"book_status"`
}

type BorrowRecord struct {
	ID              string       `json:"id"`
	BookID          string       `json:"bookId"`
	PatronID        string       `json:"patronId"`
	BorrowedAt      string       `json:"borrowedAt"`
	DueDate         string       `json:"dueDate"`
	ReturnedAt      *string      `json:"returnedAt,omitempty"`
	RenewalCount    int32        `json:"renewalCount"`
	PreviousDueDate *string      `json:"previousDueDate,omitempty"`
	Status          BorrowStatus `json:"status"`
	BookCopyID      int32        `json:"bookCopyId"`
}

func (BorrowRecord) IsRenewLoanResult() {}

type Fine struct {
	FineID            string  `json:"fine_id"`
	PatronID          string  `json:"patronId"`
	BookID            string  `json:"bookId"`
	DaysLate          int32   `json:"daysLate"`
	RatePerDay        float64 `json:"ratePerDay"`
	Amount            float64 `json:"amount"`
	CreatedAt         string  `json:"createdAt"`
	ViolationRecordID string  `json:"violationRecordId"`
}

type Mutation struct {
}

type Patron struct {
	PatronID      string        `json:"patron_id"`
	FirstName     string        `json:"first_name"`
	LastName      string        `json:"last_name"`
	PhoneNumber   string        `json:"phone_number"`
	PatronCreated string        `json:"patron_created"`
	Status        *PatronStatus `json:"status,omitempty"`
}

type PatronStatus struct {
	PatronID     string  `json:"patron_id"`
	WarningCount int32   `json:"warning_count"`
	PatronStatus Status  `json:"patron_status"`
	UnpaidFees   float64 `json:"unpaid_fees"`
}

type Query struct {
}

type RenewalError struct {
	Code    RenewalErrorCode `json:"code"`
	Message string           `json:"message"`
}

func (RenewalError) IsRenewLoanResult() {}

type Reservation struct {
	ID         string            `json:"id"`
	BookID     string            `json:"bookId"`
	PatronID   string            `json:"patronId"`
	ReservedAt string            `json:"reservedAt"`
	ExpiresAt  string            `json:"expiresAt"`
	Status     ReservationStatus `json:"status"`
}

type Subscription struct {
}

type ViolationRecord struct {
	ViolationRecordID string          `json:"violation_record_id"`
	PatronID          string          `json:"patron_id"`
	ViolationType     ViolationType   `json:"violation_type"`
	ViolationInfo     string          `json:"violation_info"`
	ViolationCreated  string          `json:"violation_created"`
	ViolationStatus   ViolationStatus `json:"violation_status"`
}

type BorrowStatus string

const (
	BorrowStatusActive   BorrowStatus = "ACTIVE"
	BorrowStatusReturned BorrowStatus = "RETURNED"
	BorrowStatusOverdue  BorrowStatus = "OVERDUE"
	BorrowStatusRenewed  BorrowStatus = "RENEWED"
)

var AllBorrowStatus = []BorrowStatus{
	BorrowStatusActive,
	BorrowStatusReturned,
	BorrowStatusOverdue,
	BorrowStatusRenewed,
}

func (e BorrowStatus) IsValid() bool {
	switch e {
	case BorrowStatusActive, BorrowStatusReturned, BorrowStatusOverdue, BorrowStatusRenewed:
		return true
	}
	return false
}

func (e BorrowStatus) String() string {
	return string(e)
}

func (e *BorrowStatus) UnmarshalGQL(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = BorrowStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid BorrowStatus", str)
	}
	return nil
}

func (e BorrowStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

func (e *BorrowStatus) UnmarshalJSON(b []byte) error {
	s, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}
	return e.UnmarshalGQL(s)
}

func (e BorrowStatus) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	e.MarshalGQL(&buf)
	return buf.Bytes(), nil
}

type RenewalErrorCode string

const (
	RenewalErrorCodeMaxRenewalsReached  RenewalErrorCode = "MAX_RENEWALS_REACHED"
	RenewalErrorCodeItemReserved        RenewalErrorCode = "ITEM_RESERVED"
	RenewalErrorCodePatronBlocked       RenewalErrorCode = "PATRON_BLOCKED"
	RenewalErrorCodeLoanNotFound        RenewalErrorCode = "LOAN_NOT_FOUND"
	RenewalErrorCodeLoanAlreadyReturned RenewalErrorCode = "LOAN_ALREADY_RETURNED"
)

var AllRenewalErrorCode = []RenewalErrorCode{
	RenewalErrorCodeMaxRenewalsReached,
	RenewalErrorCodeItemReserved,
	RenewalErrorCodePatronBlocked,
	RenewalErrorCodeLoanNotFound,
	RenewalErrorCodeLoanAlreadyReturned,
}

func (e RenewalErrorCode) IsValid() bool {
	switch e {
	case RenewalErrorCodeMaxRenewalsReached, RenewalErrorCodeItemReserved, RenewalErrorCodePatronBlocked, RenewalErrorCodeLoanNotFound, RenewalErrorCodeLoanAlreadyReturned:
		return true
	}
	return false
}

func (e RenewalErrorCode) String() string {
	return string(e)
}

func (e *RenewalErrorCode) UnmarshalGQL(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = RenewalErrorCode(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid RenewalErrorCode", str)
	}
	return nil
}

func (e RenewalErrorCode) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

func (e *RenewalErrorCode) UnmarshalJSON(b []byte) error {
	s, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}
	return e.UnmarshalGQL(s)
}

func (e RenewalErrorCode) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	e.MarshalGQL(&buf)
	return buf.Bytes(), nil
}

type ReservationStatus string

const (
	ReservationStatusPending   ReservationStatus = "PENDING"
	ReservationStatusFulfilled ReservationStatus = "FULFILLED"
	ReservationStatusCancelled ReservationStatus = "CANCELLED"
	ReservationStatusExpired   ReservationStatus = "EXPIRED"
)

var AllReservationStatus = []ReservationStatus{
	ReservationStatusPending,
	ReservationStatusFulfilled,
	ReservationStatusCancelled,
	ReservationStatusExpired,
}

func (e ReservationStatus) IsValid() bool {
	switch e {
	case ReservationStatusPending, ReservationStatusFulfilled, ReservationStatusCancelled, ReservationStatusExpired:
		return true
	}
	return false
}

func (e ReservationStatus) String() string {
	return string(e)
}

func (e *ReservationStatus) UnmarshalGQL(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ReservationStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ReservationStatus", str)
	}
	return nil
}

func (e ReservationStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

func (e *ReservationStatus) UnmarshalJSON(b []byte) error {
	s, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}
	return e.UnmarshalGQL(s)
}

func (e ReservationStatus) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	e.MarshalGQL(&buf)
	return buf.Bytes(), nil
}

type Status string

const (
	StatusGood    Status = "Good"
	StatusWarned  Status = "Warned"
	StatusBanned  Status = "Banned"
	StatusPending Status = "Pending"
)

var AllStatus = []Status{
	StatusGood,
	StatusWarned,
	StatusBanned,
	StatusPending,
}

func (e Status) IsValid() bool {
	switch e {
	case StatusGood, StatusWarned, StatusBanned, StatusPending:
		return true
	}
	return false
}

func (e Status) String() string {
	return string(e)
}

func (e *Status) UnmarshalGQL(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Status(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Status", str)
	}
	return nil
}

func (e Status) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

func (e *Status) UnmarshalJSON(b []byte) error {
	s, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}
	return e.UnmarshalGQL(s)
}

func (e Status) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	e.MarshalGQL(&buf)
	return buf.Bytes(), nil
}

type ViolationStatus string

const (
	ViolationStatusOngoing  ViolationStatus = "Ongoing"
	ViolationStatusResolved ViolationStatus = "Resolved"
)

var AllViolationStatus = []ViolationStatus{
	ViolationStatusOngoing,
	ViolationStatusResolved,
}

func (e ViolationStatus) IsValid() bool {
	switch e {
	case ViolationStatusOngoing, ViolationStatusResolved:
		return true
	}
	return false
}

func (e ViolationStatus) String() string {
	return string(e)
}

func (e *ViolationStatus) UnmarshalGQL(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ViolationStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ViolationStatus", str)
	}
	return nil
}

func (e ViolationStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

func (e *ViolationStatus) UnmarshalJSON(b []byte) error {
	s, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}
	return e.UnmarshalGQL(s)
}

func (e ViolationStatus) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	e.MarshalGQL(&buf)
	return buf.Bytes(), nil
}

type ViolationType string

const (
	ViolationTypeLateReturn  ViolationType = "Late_Return"
	ViolationTypeUnpaidFees  ViolationType = "Unpaid_Fees"
	ViolationTypeDamagedBook ViolationType = "Damaged_Book"
)

var AllViolationType = []ViolationType{
	ViolationTypeLateReturn,
	ViolationTypeUnpaidFees,
	ViolationTypeDamagedBook,
}

func (e ViolationType) IsValid() bool {
	switch e {
	case ViolationTypeLateReturn, ViolationTypeUnpaidFees, ViolationTypeDamagedBook:
		return true
	}
	return false
}

func (e ViolationType) String() string {
	return string(e)
}

func (e *ViolationType) UnmarshalGQL(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ViolationType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ViolationType", str)
	}
	return nil
}

func (e ViolationType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

func (e *ViolationType) UnmarshalJSON(b []byte) error {
	s, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}
	return e.UnmarshalGQL(s)
}

func (e ViolationType) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	e.MarshalGQL(&buf)
	return buf.Bytes(), nil
}
