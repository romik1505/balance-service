package handler

import (
	"errors"
	"log"
	"net/http"
)

var (
	ErrorBadRequest = errors.New("bad request")
)

func ErrorWrapper(fn func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := fn(w, r)
		if err != nil {
			log.Printf("Error: %v", err)
			if errors.Is(err, ErrorBadRequest) {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
