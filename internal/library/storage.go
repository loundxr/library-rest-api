package library

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type Storage struct {
	uniqueness map[string]struct{}
	sync.RWMutex
	lib map[uuid.UUID]Book
}

func NewStorage() *Storage {
	return &Storage{
		uniqueness: make(map[string]struct{}),
		lib:        make(map[uuid.UUID]Book),
	}
}

func (s *Storage) AddBook(params BookParams) (*Book, error) {
	newParams := fmt.Sprintf("%s/%s/%d", params.Title, params.Author, params.Pages)

	s.Lock()
	defer s.Unlock()

	if _, ok := s.uniqueness[newParams]; ok {
		return nil, ErrBookAlreadyExists
	}
	s.uniqueness[newParams] = struct{}{}

	book := NewBook(params.Title, params.Author, params.Pages)
	s.lib[book.ID] = *book
	return book, nil
}
