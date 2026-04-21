package library

import "errors"

var ErrBookAlreadyExists = errors.New("book already exists")
var ErrBookNotFound = errors.New("book not found")
var ErrInvalidSortKey = errors.New("invalid sort key in query parameters")
