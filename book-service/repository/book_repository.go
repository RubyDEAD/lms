package repository

import (
	"context"
	"log"
	"book-service/models"

	"github.com/jackc/pgx/v5"
)

type BookRepository struct {
	DB *pgx.Conn
}

func NewBookRepository(conn *pgx.Conn) *BookRepository {
	return &BookRepository{DB: conn}
}

func (repo *BookRepository) GetBooks() ([]models.Book, error) {
	rows, err := repo.DB.Query(context.Background(), "SELECT id, title, author, published_year FROM books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.PublishedYear)
		if err != nil {
			log.Println(err)
			continue
		}
		books = append(books, book)
	}
	return books, nil
}

func (repo *BookRepository) AddBook(title, author string, year int) (models.Book, error) {
	var book models.Book
	err := repo.DB.QueryRow(
		context.Background(),
		"INSERT INTO books (title, author, published_year) VALUES ($1, $2, $3) RETURNING id, title, author, published_year",
		title, author, year,
	).Scan(&book.ID, &book.Title, &book.Author, &book.PublishedYear)

	return book, err
}
