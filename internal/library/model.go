package library

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Book struct {
	ID             uuid.UUID   `json:"id"`
	Title          string      `json:"title"`
	Author         string      `json:"author"`
	Pages          int         `json:"pages"`
	TimeOfAddition PrettyTime  `json:"time_of_addition"`
	IsRead         bool        `json:"is_read"`
	ReadAt         *PrettyTime `json:"read_at"`
}

func (b *Book) Read() {
	b.IsRead = true
	prettyNow := PrettyTime(time.Now())
	b.ReadAt = &prettyNow
}

func (b *Book) Unread() {
	b.IsRead = false
	b.ReadAt = nil
}

type BookParams struct {
	Title  string
	Author string
	Pages  int
}

type GetBooksParams struct {
	Author   string
	IsRead   *bool
	SortType string
}

func NewBook(title, author string, pages int) *Book {
	return &Book{
		ID:             uuid.New(),
		Title:          title,
		Author:         author,
		Pages:          pages,
		TimeOfAddition: PrettyTime(time.Now()),
		IsRead:         false,
		ReadAt:         nil,
	}
}

type PrettyTime time.Time

func (pt PrettyTime) MarshalJSON() ([]byte, error) {
	t := time.Time(pt)
	str := fmt.Sprintf("\"%s\"", t.Format("02.01.2006 15:04:05"))
	return []byte(str), nil
}
