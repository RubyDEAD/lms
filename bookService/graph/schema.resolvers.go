package graph

import (
	"context"
	"fmt"
	"time"

	"github.com/Cat6utpcableclarke/bookService/graph/model"
)

// AddBook adds a new book along with the author (if the author does not exist).
func (r *mutationResolver) AddBook(ctx context.Context, title string, authorName string, datePublished string, description string) (*model.Book, error) {
	var authorID int
	err := r.DB.QueryRow(ctx, "INSERT INTO authors (author_name) VALUES ($1) ON CONFLICT (author_name) DO UPDATE SET author_name=EXCLUDED.author_name RETURNING id", authorName).Scan(&authorID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert or get author: %v", err)
	}

	var bookID string
	err = r.DB.QueryRow(ctx, "INSERT INTO books (title, author_id, date_published, description) VALUES ($1, $2, $3, $4) RETURNING id", title, authorID, datePublished, description).Scan(&bookID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert book: %v", err)
	}

	if bookID != "" {
		_, err = r.DB.Exec(ctx, "INSERT INTO book_copies (book_id) VALUES ($1)", bookID)
		if err != nil {
			return nil, fmt.Errorf("failed to insert book copy: %v", err)
		}

	}

	return &model.Book{
		ID:            bookID,
		Title:         title,
		AuthorName:    authorName,
		DatePublished: datePublished,
		Description:   description,
	}, nil
}

// UpdateBook updates an existing book's details.
func (r *mutationResolver) UpdateBook(ctx context.Context, id string, title *string, authorName string, datePublished *string, description *string) (*model.Book, error) {
	var authorID int
	err := r.DB.QueryRow(ctx, "INSERT INTO authors (author_name) VALUES ($1) ON CONFLICT (author_name) DO UPDATE SET author_name=EXCLUDED.author_name RETURNING id", authorName).Scan(&authorID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert or get author: %v", err)
	}

	_, err = r.DB.Exec(ctx, "UPDATE books SET title=COALESCE($1, title), author_id=$2, date_published=COALESCE($3, date_published), description=COALESCE($4, description) WHERE id=$5", title, authorID, datePublished, description, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update book: %v", err)
	}

	return &model.Book{
		ID: id,
		Title: func() string {
			if title != nil {
				return *title
			} else {
				return ""
			}
		}(),
		AuthorName: authorName,
		DatePublished: func() string {
			if datePublished != nil {
				return *datePublished
			}
			return ""
		}(),
		Description: func() string {
			if description != nil {
				return *description
			}
			return ""
		}(),
	}, nil
}

// UpdateBookCopyStatus updates the status of a book copy.
func (r *mutationResolver) UpdateBookCopyStatus(ctx context.Context, id string, bookStatus *string) (*model.BookCopies, error) {
	_, err := r.DB.Exec(ctx, "UPDATE book_copies SET book_status=$1 WHERE id=$2", bookStatus, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update book copy status: %v", err)
	}

	return &model.BookCopies{
		ID: id,
		BookStatus: func() string {
			if bookStatus != nil {
				return *bookStatus
			}
			return ""
		}(),
	}, nil
}

// DeleteBook deletes a book by ID.
func (r *mutationResolver) DeleteBook(ctx context.Context, id string) (bool, error) {
	_, err := r.DB.Exec(ctx, "DELETE FROM books WHERE id=$1", id)
	if err != nil {
		return false, fmt.Errorf("failed to delete book: %v", err)
	}
	return true, nil
}

// GetFilteredBooks fetches books based on filters.
func (r *queryResolver) GetFilteredBooks(ctx context.Context, filter *model.BookFilter) ([]*model.Book, error) {
	return nil, fmt.Errorf("not implemented: GetFilteredBooks")
}

// GetBookByID fetches a book by its ID.
func (r *queryResolver) GetBookByID(ctx context.Context, id string) (*model.Book, error) {
	var book model.Book
	var authorName string
	var datePublished time.Time // Store as time.Time to handle DATE type

	err := r.DB.QueryRow(ctx, "SELECT b.id, b.title, a.author_name, b.date_published, b.description FROM books b JOIN authors a ON b.author_id = a.id WHERE b.id=$1", id).
		Scan(&book.ID, &book.Title, &authorName, &datePublished, &book.Description)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch book: %v", err)
	}

	// Convert `datePublished` to string (ISO 8601 format)
	formattedDate := datePublished.Format("2006-01-02")
	book.DatePublished = formattedDate
	book.AuthorName = authorName

	return &book, nil
}

// GetBookCopiesByID fetches all copies of a book by book ID.
func (r *queryResolver) GetBookCopiesByID(ctx context.Context, id string) ([]*model.BookCopies, error) {
	rows, err := r.DB.Query(ctx, "SELECT id, book_id, book_status, FROM book_copies WHERE book_id=$1", id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch book copies: %v", err)
	}
	defer rows.Close()

	var copies []*model.BookCopies
	for rows.Next() {
		var copy model.BookCopies
		err := rows.Scan(&copy.ID, &copy.BookID, &copy.BookStatus)
		if err != nil {
			return nil, fmt.Errorf("failed to scan book copy: %v", err)
		}
		copies = append(copies, &copy)
	}
	return copies, nil
}

// SearchBooks allows searching books by title or author name.
func (r *queryResolver) SearchBooks(ctx context.Context, query string) ([]*model.Book, error) {
	rows, err := r.DB.Query(ctx, "SELECT b.id, b.title, a.author_name, b.date_published, b.description FROM books b JOIN authors a ON b.author_id = a.id WHERE b.title ILIKE $1 OR a.author_name ILIKE $1", "%"+query+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to search books: %v", err)
	}
	defer rows.Close()

	var books []*model.Book
	for rows.Next() {
		var book model.Book
		var authorName string
		err := rows.Scan(&book.ID, &book.Title, &authorName, &book.DatePublished, &book.Description)
		if err != nil {
			return nil, fmt.Errorf("failed to scan book: %v", err)
		}
		book.AuthorName = authorName
		books = append(books, &book)
	}
	return books, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
