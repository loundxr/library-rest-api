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
	lib library.BookRepository
}

func NewHTTPHandlers(lib library.BookRepository) *HTTPHandlers {
	return &HTTPHandlers{
		lib: lib,
	}
}

func parseID(str string, w http.ResponseWriter) (uuid.UUID, error) {
	id, err := uuid.Parse(str)
	if err != nil {
		SendError(w, err.Error(), http.StatusBadRequest)
		return uuid.UUID{}, err
	}
	return id, nil
}

func (h *HTTPHandlers) HandleAddBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
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
		Year:   dto.Year,
		Review: dto.Review,
	}

	book, err := h.lib.AddBook(ctx, params)
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
	ctx := r.Context()
	query := r.URL.Query()

	params := library.GetBooksParams{
		Author: query.Get("author"),
	}

	if query.Has("is_read") {
		isRead, err := strconv.ParseBool(query.Get("is_read"))
		if err != nil {
			SendError(w, "invalid is_read parameter", http.StatusBadRequest)
			return
		}
		params.IsRead = &isRead
	}

	limit := 10
	if query.Has("limit") {
		v, err := strconv.Atoi(query.Get("limit"))
		if err != nil {
			SendError(w, err.Error(), http.StatusBadRequest)
			return
		}
		if v > 0 {
			if v > 100 {
				v = 100
			}
			limit = v
		}
	}
	params.Limit = &limit

	offset := 0
	if query.Has("offset") {
		v, err := strconv.Atoi(query.Get("offset"))
		if err != nil {
			SendError(w, err.Error(), http.StatusBadRequest)
			return
		}
		if v > 0 {
			offset = v
		}
	}
	params.Offset = &offset

	if query.Has("year") {
		year, err := strconv.Atoi(query.Get("year"))
		if err != nil {
			SendError(w, err.Error(), http.StatusBadRequest)
			return
		}
		params.Year = &year
	}

	params.SortType = query.Get("sort")

	books, err := h.lib.GetAllBooks(ctx, params)

	if err != nil {
		if errors.Is(err, library.ErrInvalidSortKey) {
			SendError(w, err.Error(), http.StatusBadRequest)
		} else {
			SendError(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(books); err != nil {
		fmt.Println("error encoding json:", err)
		return
	}
}

func (h *HTTPHandlers) HandleGetBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := parseID(mux.Vars(r)["id"], w)
	if err != nil {
		return
	}

	book, err := h.lib.GetBook(ctx, id)
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
	ctx := r.Context()
	id, err := parseID(mux.Vars(r)["id"], w)
	if err != nil {
		return
	}

	if err := h.lib.DeleteBook(ctx, id); err != nil {
		if errors.Is(err, library.ErrBookNotFound) {
			SendError(w, err.Error(), http.StatusNotFound)
		} else {
			SendError(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *HTTPHandlers) HandlePatchBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := parseID(mux.Vars(r)["id"], w)
	if err != nil {
		return
	}

	var dto PatchDTO

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	params := library.UpdateBookParams{
		Title:  dto.Title,
		Author: dto.Author,
		Pages:  dto.Pages,
		IsRead: dto.IsRead,
		Year:   dto.Year,
		Review: dto.Review,
	}

	b, err := h.lib.PatchBook(ctx, id, params)
	if err != nil {
		if errors.Is(err, library.ErrBookNotFound) {
			SendError(w, err.Error(), http.StatusNotFound)
		} else if errors.Is(err, library.ErrBookAlreadyExists) {
			SendError(w, err.Error(), http.StatusConflict)
		} else {
			SendError(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(b); err != nil {
		fmt.Println("error encoding json:", err)
	}
}
