package graph

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

// Helper function: Call Borrowing Service to get Days Late
func getDaysLateFromBorrowingService(patronID string, bookID string) (int, error) {
	query := `
		query GetBorrowRecord($patronId: ID!, $bookId: ID!) {
			borrowRecordByPatronAndBook(patronId: $patronId, bookId: $bookId) {
				dueDate
				returnedAt
			}
		}
	`
	requestBody, _ := json.Marshal(map[string]interface{}{
		"query": query,
		"variables": map[string]interface{}{
			"patronId": patronID,
			"bookId": bookID,
		},
	})

	resp, err := http.Post("http://localhost:8082/query", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			BorrowRecordByPatronAndBook struct {
				DueDate    string `json:"dueDate"`
				ReturnedAt string `json:"returnedAt"`
			} `json:"borrowRecordByPatronAndBook"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	layout := "2006-01-02"
	dueDate, err1 := time.Parse(layout, result.Data.BorrowRecordByPatronAndBook.DueDate)
	returnedAt, err2 := time.Parse(layout, result.Data.BorrowRecordByPatronAndBook.ReturnedAt)

	if err1 != nil || err2 != nil {
		return 0, err
	}

	daysLate := int(returnedAt.Sub(dueDate).Hours() / 24)
	if daysLate < 0 {
		daysLate = 0
	}

	return daysLate, nil
}
