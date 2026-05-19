package library

import (
	"context"

	"github.com/google/uuid"
)

type BookRepository interface {
	AddBook(ctx context.Context, params BookParams) (Book, error)
	GetAllBooks(ctx context.Context, p GetBooksParams) ([]Book, error)
	GetBook(ctx context.Context, id uuid.UUID) (Book, error)
	DeleteBook(ctx context.Context, id uuid.UUID) error
	MarkBook(ctx context.Context, id uuid.UUID, read bool) (Book, error)
	PatchBook(ctx context.Context, id uuid.UUID, p UpdateBookParams) (Book, error)
}
