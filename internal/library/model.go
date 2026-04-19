package library

import (
	"time"

	"github.com/google/uuid"
)

type Book struct {
	Id             uuid.UUID
	Title          string
	Author         string
	Pages          int
	TimeOfAddition time.Time
	IsRead         bool
	ReadAt         *time.Time
}
