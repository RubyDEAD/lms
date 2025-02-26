package graph

import "book-service/graph/model"

// Resolver serves as dependency injection for your app.
type Resolver struct {
	Books []*model.Book
}
