package api

import "github.com/gorilla/mux"

type HTTPServer struct {
	httpHandlers *HTTPHandlers
}

func NewHTTPServer(handlers *HTTPHandlers) *HTTPServer {
	return &HTTPServer{
		httpHandlers: handlers,
	}
}

func StartServer() {
	r := mux.NewRouter()

}
