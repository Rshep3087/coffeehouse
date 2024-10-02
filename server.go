package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rshep3087/coffeehouse/cache"
	"github.com/rshep3087/coffeehouse/postgres"
	"go.uber.org/zap"
)

//go:generate moq -out pubsub_moq_test.go . PubSub
type PubSub interface {
	// Publish publishes a message to a topic without waiting for a response
	Publish(topic string, data []byte) error
}

type server struct {
	router  *httprouter.Router
	log     *zap.SugaredLogger
	queries *postgres.Queries
	pubsub  PubSub
	cacher  cache.RecipeCacher
}

func newServer(ps PubSub, cacher cache.RecipeCacher) *server {
	s := &server{
		router: httprouter.New(),
		pubsub: ps,
		cacher: cacher,
	}
	s.routes()
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
