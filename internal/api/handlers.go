package api

import (
	"encoding/json"
	"errors"
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
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := dto.Validate(); err != nil {
		SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	params := library.BookParams{
		Title:  dto.Title,
		Author: dto.Author,
		Pages:  dto.Pages,
	}

	book, err := h.lib.AddBook(params)
	if err != nil {
		if errors.Is(err, library.ErrBookAlreadyExists) {
			SendError(w, err.Error(), http.StatusConflict)
		} else {
			SendError(w, err.Error(), http.StatusInternalServerError)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(book); err != nil {
		SendError(w, err.Error(), http.StatusInternalServerError)
	}
}
