package library

import "errors"

var ErrBookAlreadyExists = errors.New("book already exists")
var ErrBookNotFound = errors.New("book not found")
