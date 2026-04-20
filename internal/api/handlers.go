package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"library-rest-api/internal/library"
	"net/http"
	"strconv"
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
		SendError(w, err.Error(), http.StatusBadRequest)
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
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(book); err != nil {
		fmt.Println("error encoding json:", err)
	}
}

func (h *HTTPHandlers) HandleGetAllBooks(w http.ResponseWriter, r *http.Request) {
	author := r.URL.Query().Get("author")
	params := library.GetBooksParams{
		Author: author,
	}

	if r.URL.Query().Has("is_read") {
		isReadStr := r.URL.Query().Get("is_read")

		isRead, err := strconv.ParseBool(isReadStr)
		if err != nil {
			SendError(w, err.Error(), http.StatusBadRequest)
			return
		}
		params.IsRead = &isRead
	}

	books := h.lib.GetAllBooks(params)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(books); err != nil {
		fmt.Println("error encoding json:", err)
		return
	}
}
