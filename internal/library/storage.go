package library

import (
	"sync"

	"github.com/google/uuid"
)

type Storage struct {
	sync.RWMutex
	Lib map[uuid.UUID]string
}
