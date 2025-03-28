// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type Membership struct {
	MembershipID string          `json:"membership_id"`
	PatronID     string          `json:"patron_id"`
	Level        MembershipLevel `json:"level"`
}

type Mutation struct {
}

type Patron struct {
	PatronID    string             `json:"patron_id"`
	FirstName   string             `json:"first_name"`
	LastName    string             `json:"last_name"`
	PhoneNumber string             `json:"phone_number"`
	Membership  *Membership        `json:"membership,omitempty"`
	Status      *PatronStatus      `json:"status,omitempty"`
	Violations  []*ViolationRecord `json:"violations,omitempty"`
}

type PatronStatus struct {
	PatronID     string  `json:"patron_id"`
	WarningCount int32   `json:"warning_count"`
	PatronStatus Status  `json:"patron_status"`
	UnpaidFees   float64 `json:"unpaid_fees"`
}

type Query struct {
}

type ViolationRecord struct {
	ViolationRecordID string        `json:"violation_record_id"`
	PatronID          string        `json:"patron_id"`
	ViolationType     ViolationType `json:"violation_type"`
	ViolationInfo     string        `json:"violation_info"`
}

type MembershipLevel string

const (
	MembershipLevelBronze MembershipLevel = "Bronze"
	MembershipLevelSilver MembershipLevel = "Silver"
	MembershipLevelGold   MembershipLevel = "Gold"
)

var AllMembershipLevel = []MembershipLevel{
	MembershipLevelBronze,
	MembershipLevelSilver,
	MembershipLevelGold,
}

func (e MembershipLevel) IsValid() bool {
	switch e {
	case MembershipLevelBronze, MembershipLevelSilver, MembershipLevelGold:
		return true
	}
	return false
}

func (e MembershipLevel) String() string {
	return string(e)
}

func (e *MembershipLevel) UnmarshalGQL(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = MembershipLevel(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid MembershipLevel", str)
	}
	return nil
}

func (e MembershipLevel) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
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
