package api

import (
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type HTTPServer struct {
	httpHandlers *HTTPHandlers
}

func NewHTTPServer(handlers *HTTPHandlers) *HTTPServer {
	return &HTTPServer{
		httpHandlers: handlers,
	}
}

func (s *HTTPServer) Start(addr string) error {
	r := mux.NewRouter()

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
