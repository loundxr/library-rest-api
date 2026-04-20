package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"library-rest-api/internal/library"
	"net/http"
	"time"
)

type BookDTO struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Pages  int    `json:"pages"`
}

func (b BookDTO) Validate() error {
	if b.Title == "" {
		return errors.New("empty field: title")
	}

	if b.Author == "" {
		return errors.New("empty field: author")
	}

	if b.Pages <= 0 {
		return errors.New("number of pages must be a positive number")
	}

	return nil
}

type ErrorDTO struct {
	Message string
	Time    library.PrettyTime
}

func SendError(w http.ResponseWriter, msg string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	e := ErrorDTO{
		Message: msg,
		Time:    library.PrettyTime(time.Now()),
	}
	if err := json.NewEncoder(w).Encode(e); err != nil {
		fmt.Println("error encoding errorDTO:", err)
	}
}

func (e ErrorDTO) String() string {
	defer func() {
		if p := recover(); p != nil {
			fmt.Println("recovered from panic:", p)
		}
	}()
	b, err := json.MarshalIndent(e, "", "    ")
	if err != nil {
		panic(err)
	}
	return string(b)
}
