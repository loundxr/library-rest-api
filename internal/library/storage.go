package library

import (
	"fmt"
	"strings"
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

func GenerateKey(title, author string, pages int) string {
	return fmt.Sprintf("%s/%s/%d",
		strings.ToLower(strings.TrimSpace(title)),
		strings.ToLower(strings.TrimSpace(author)),
		pages)
}

func (s *Storage) AddBook(params BookParams) (*Book, error) {
	newParams := GenerateKey(params.Title, params.Author, params.Pages)

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

func (s *Storage) GetAllBooks(p GetBooksParams) []Book {
	s.RLock()
	defer s.RUnlock()

	p.Author = strings.ToLower(strings.TrimSpace(p.Author))

	cpy := make([]Book, 0)

	for _, b := range s.lib {
		if p.Author != "" {
			bAuthor := strings.ToLower(strings.TrimSpace(b.Author))
			if !strings.Contains(bAuthor, p.Author) {
				continue
			}
		}

		if p.IsRead != nil && *p.IsRead != b.IsRead {
			continue
		}
		cpy = append(cpy, b)
	}

	return cpy
}

func (s *Storage) GetBook(id uuid.UUID) (*Book, error) {
	s.RLock()
	defer s.RUnlock()

	book, ok := s.lib[id]
	if !ok {
		return nil, ErrBookNotFound
	}
	return &book, nil
}

func (s *Storage) DeleteBook(id uuid.UUID) error {
	s.Lock()
	defer s.Unlock()

	b, ok := s.lib[id]
	if !ok {
		return ErrBookNotFound
	}

	key := GenerateKey(b.Title, b.Author, b.Pages)

	delete(s.lib, id)
	delete(s.uniqueness, key)
	return nil
}

func (s *Storage) MarkBook(id uuid.UUID, read bool) (Book, error) {
	s.Lock()
	defer s.Unlock()

	book, ok := s.lib[id]
	if !ok {
		return Book{}, ErrBookNotFound
	}

	if read {
		book.Read()
	} else {
		book.Unread()
	}
	s.lib[id] = book
	return book, nil
}
