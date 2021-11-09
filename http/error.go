package http

import (
	"fmt"
	"net/http"
)

type appError struct {
	Error   error
	Message string
	Code    int
}

type ErrHandlerFunc func(w http.ResponseWriter, r *http.Request) *appError

func (s *Server) Error(hf ErrHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if e := hf(w, r); e != nil {
			if e.Message == "" {
				e.Message = "An error has occurred."
			}
			fmt.Printf("%+v\n", e.Error)
			http.Error(w, e.Message, e.Code)
		}
	}
}
