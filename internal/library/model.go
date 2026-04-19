package library

import (
	"time"

	"github.com/google/uuid"
)

type Book struct {
	ID             uuid.UUID  `json:"id"`
	Title          string     `json:"title"`
	Author         string     `json:"author"`
	Pages          int        `json:"pages"`
	TimeOfAddition time.Time  `json:"time_of_addition"`
	IsRead         bool       `json:"is_read"`
	ReadAt         *time.Time `json:"read_at"`
}

type BookParams struct {
	Title  string
	Author string
	Pages  int
}

func NewBook(title, author string, pages int) *Book {
	return &Book{
		ID:             uuid.New(),
		Title:          title,
		Author:         author,
		Pages:          pages,
		TimeOfAddition: time.Now(),
		IsRead:         false,
		ReadAt:         nil,
	}
}
