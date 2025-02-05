package web

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
	router              *httprouter.Router
	log                 *zap.SugaredLogger
	Queries             *postgres.Queries
	pubsub              PubSub
	cacher              cache.RecipeCacher
	coffeeImageProvider CoffeeImageProvider
}

func NewServer(
	log *zap.SugaredLogger,
	ps PubSub,
	cacher cache.RecipeCacher,
	coffeeImageProvider CoffeeImageProvider,
) *server {
	router := httprouter.New()
	router.PanicHandler = func(w http.ResponseWriter, r *http.Request, i interface{}) {
		log.Errorw("recovered from panic", "error", i)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	s := &server{
		log:                 log,
		router:              httprouter.New(),
		pubsub:              ps,
		cacher:              cacher,
		coffeeImageProvider: coffeeImageProvider,
	}
	s.routes()
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
