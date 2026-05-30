package api

import (
	"errors"
	"log"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

type HTTPServer struct {
	httpHandlers *HTTPHandlers
	logger       *slog.Logger
}

func NewHTTPServer(handlers *HTTPHandlers, logger *slog.Logger) *HTTPServer {
	return &HTTPServer{
		httpHandlers: handlers,
		logger:       logger,
	}
}

func (s *HTTPServer) Start(addr string) error {
	r := mux.NewRouter()

	r.Use(s.LoggingMiddleware)

	r.HandleFunc("/books", s.httpHandlers.HandleAddBook).Methods(http.MethodPost)
	r.HandleFunc("/books", s.httpHandlers.HandleGetAllBooks).Methods(http.MethodGet)
	r.HandleFunc("/books/{id}", s.httpHandlers.HandleGetBook).Methods(http.MethodGet)
	r.HandleFunc("/books/{id}", s.httpHandlers.HandleDeleteBook).Methods(http.MethodDelete)
	r.HandleFunc("/books/{id}", s.httpHandlers.HandlePatchBook).Methods(http.MethodPatch)

	log.Println("starting server on:", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		return err
	}
	return nil
}
