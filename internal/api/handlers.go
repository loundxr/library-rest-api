package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"library-rest-api/internal/library"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type HTTPHandlers struct {
	lib *library.Storage
}

func NewHTTPHandlers(lib *library.Storage) *HTTPHandlers {
	return &HTTPHandlers{
		lib: lib,
	}
}

func ParseID(str string, w http.ResponseWriter) (uuid.UUID, error) {
	id, err := uuid.Parse(str)
	if err != nil {
		SendError(w, err.Error(), http.StatusBadRequest)
		return uuid.UUID{}, err
	}
	return id, nil
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
			SendError(w, "internal server error", http.StatusInternalServerError)
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

func (h *HTTPHandlers) HandleGetBook(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(mux.Vars(r)["id"], w)
	if err != nil {
		return
	}

	book, err := h.lib.GetBook(id)
	if err != nil {
		if errors.Is(err, library.ErrBookNotFound) {
			SendError(w, err.Error(), http.StatusNotFound)
		} else {
			SendError(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(book); err != nil {
		fmt.Println("error encoding json:", err)
		return
	}
}

func (h *HTTPHandlers) HandleDeleteBook(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(mux.Vars(r)["id"], w)
	if err != nil {
		return
	}

	if err := h.lib.DeleteBook(id); err != nil {
		if errors.Is(err, library.ErrBookNotFound) {
			SendError(w, err.Error(), http.StatusNotFound)
		} else {
			SendError(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *HTTPHandlers) HandleMarkReadBook(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(mux.Vars(r)["id"], w)
	if err != nil {
		return
	}

	var readDTO ReadDTO

	if err := json.NewDecoder(r.Body).Decode(&readDTO); err != nil {
		SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	b, err := h.lib.MarkBook(id, readDTO.Read)

	if err != nil {
		if errors.Is(err, library.ErrBookNotFound) {
			SendError(w, err.Error(), http.StatusNotFound)
		} else {
			SendError(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(b); err != nil {
		fmt.Println("error encoding json:", err)
	}

}
