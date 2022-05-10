package main

import (
	"fmt"
	"net/http"
)

func (s *server) routes() {
	s.router.HandlerFunc(http.MethodGet, "/health", s.health())
}

func (s *server) health() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "healthy")
	}
}
