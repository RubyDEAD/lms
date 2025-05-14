package graph

import (
	"sync"

	"github.com/Cat6utpcableclarke/bookService/graph/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB                 *pgxpool.Pool
	BookAddedObservers map[string]chan *model.Book
	mu                 sync.Mutex
}

func NewResolver(db *pgxpool.Pool) *Resolver {
	return &Resolver{
		DB:                 db,
		BookAddedObservers: make(map[string]chan *model.Book),
	}
}
