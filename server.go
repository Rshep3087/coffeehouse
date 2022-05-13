package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rshep3087/coffeehouse/postgres"
	"go.uber.org/zap"
)

type server struct {
	router  *httprouter.Router
	log     *zap.SugaredLogger
	queries *postgres.Queries
}

func newServer() *server {
	s := &server{
		router: httprouter.New(),
	}
	s.routes()
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
