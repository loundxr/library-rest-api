package api

import (
	"library-rest-api/internal/library"
	"net/http"
)

type HTTPHandlers struct {
	lib *library.Storage
}

func NewHTTPHandlers(lib *library.Storage) *HTTPHandlers {
	return &HTTPHandlers{
		lib: lib,
	}
}

func (h *HTTPHandlers) HandleAddBook(w http.ResponseWriter, r *http.Request) {
	var dto BookDTO

}
