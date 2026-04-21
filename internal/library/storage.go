package library

import (
	"cmp"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Storage struct {
	uniqueness map[string]struct{}
	mu         sync.RWMutex
	lib        map[uuid.UUID]Book
}

func NewStorage() *Storage {
	return &Storage{
		uniqueness: make(map[string]struct{}),
		lib:        make(map[uuid.UUID]Book),
	}
}

func generateKey(title, author string, pages int) string {
	return fmt.Sprintf("%s/%s/%d",
		strings.ToLower(strings.TrimSpace(title)),
		strings.ToLower(strings.TrimSpace(author)),
		pages)
}

func (s *Storage) AddBook(params BookParams) (Book, error) {
	newParams := generateKey(params.Title, params.Author, params.Pages)

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.uniqueness[newParams]; ok {
		return Book{}, ErrBookAlreadyExists
	}
	s.uniqueness[newParams] = struct{}{}

	book := NewBook(params.Title, params.Author, params.Pages)
	s.lib[book.ID] = *book
	return *book, nil
}

func (s *Storage) GetAllBooks(p GetBooksParams) ([]Book, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	p.Author = strings.ToLower(strings.TrimSpace(p.Author))

	res := make([]Book, 0)

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
		res = append(res, b)
	}

	switch p.SortType {
	case "author":
		slices.SortFunc(res, func(a, b Book) int {
			return strings.Compare(
				strings.ToLower(a.Author),
				strings.ToLower(b.Author),
			)
		})
	case "title":
		slices.SortFunc(res, func(a, b Book) int {
			return strings.Compare(
				strings.ToLower(a.Title),
				strings.ToLower(b.Title),
			)
		})
	case "pages":
		slices.SortFunc(res, func(a, b Book) int {
			return cmp.Compare(a.Pages, b.Pages)
		})
	case "time":
		slices.SortFunc(res, func(a, b Book) int {
			return time.Time(a.TimeOfAddition).Compare(time.Time(b.TimeOfAddition))
		})
	case "":
		break
	default:
		return nil, ErrInvalidSortKey
	}
	return res, nil
}

func (s *Storage) GetBook(id uuid.UUID) (Book, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	book, ok := s.lib[id]
	if !ok {
		return Book{}, ErrBookNotFound
	}
	return book, nil
}

func (s *Storage) DeleteBook(id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	b, ok := s.lib[id]
	if !ok {
		return ErrBookNotFound
	}

	key := generateKey(b.Title, b.Author, b.Pages)

	delete(s.lib, id)
	delete(s.uniqueness, key)
	return nil
}

func (s *Storage) MarkBook(id uuid.UUID, read bool) (Book, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

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
