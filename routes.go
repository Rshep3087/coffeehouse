package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/rshep3087/coffeehouse/postgres"
)

func (s *server) routes() {
	s.router.HandlerFunc(http.MethodGet, "/health", s.health())
	s.router.HandlerFunc(http.MethodPost, "/v1/recipes", s.handleCreateRecipe())
	s.router.HandlerFunc(http.MethodGet, "/v1/recipes/:id", s.handleGetRecipe())
	s.router.HandlerFunc(http.MethodGet, "/v1/recipes", s.handleGetRecipes())
}

func (s *server) health() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "healthy")
	}
}

func (s *server) handleCreateRecipe() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		recip, err := decode[postgres.CreateRecipeParams](r)
		if err != nil {
			http.Error(w, "error decoding recipe", http.StatusBadRequest)
			return
		}

		log := s.log.With(
			"method", recip.BrewMethod,
			"recipe_name", recip.RecipeName,
		)

		log.Info("adding recipe")

		newRecipe, err := s.queries.CreateRecipe(r.Context(), recip)
		if err != nil {
			log.Error(err)
			http.Error(w, "error adding recipe", http.StatusInternalServerError)
			return
		}

		// publish event that a new recipe was added
		recipeJSON, err := json.Marshal(newRecipe)
		if err != nil {
			log.Error(err)
			http.Error(w, "error marshaling new recipe", http.StatusInternalServerError)
			return
		}

		if err := s.pubsub.Publish("recipe.new", recipeJSON); err != nil {
			log.Error(err)
			http.Error(w, "error publishing new recipe event", http.StatusInternalServerError)
			return
		}

		log.Info("recipe added")
		if err = encode(w, http.StatusCreated, newRecipe); err != nil {
			log.Error(err)
			http.Error(w, "error encoding response", http.StatusInternalServerError)
			return
		}
	}
}

func (s *server) handleGetRecipe() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := httprouter.ParamsFromContext(r.Context())
		recipeID := params.ByName("id")

		log := s.log.With("id", recipeID)

		rid, err := strconv.Atoi(recipeID)
		if err != nil {
			log.Error(err)
			http.Error(w, "invalid id param", http.StatusBadRequest)
			return
		}

		log.Info("checking cache")
		cachedRecipe, err := s.cacher.GetRecipe(r.Context(), int64(rid))
		if err == nil && cachedRecipe != nil {
			log.Info("cache hit")
			if err = encode(w, http.StatusOK, cachedRecipe); err != nil {
				log.Error(err)
				http.Error(w, "error encoding response", http.StatusInternalServerError)
				return
			}
			return
		}

		log.Info("getting recipe")
		recipe, err := s.queries.GetRecipe(r.Context(), int64(rid))
		if err != nil {
			log.Error(err)
			http.Error(w, "error getting recipe", http.StatusInternalServerError)
			return
		}

		log.Info("recipe got", "recipe_name", recipe.RecipeName)

		log.Info("setting cache")
		if err := s.cacher.SetRecipe(r.Context(), int64(rid), &recipe); err != nil {
			log.Error(err)
		}

		if err = encode(w, http.StatusOK, recipe); err != nil {
			log.Error(err)
			http.Error(w, "error encoding response", http.StatusInternalServerError)
			return
		}
	}
}

func (s *server) handleGetRecipes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.log.Info("getting recipes")

		recipes, err := s.queries.ListRecipes(r.Context())
		if err != nil {
			s.log.Error(err)
			http.Error(w, "error getting recipes", http.StatusInternalServerError)
			return
		}

		s.log.Info("recipes got")

		if err := encode(w, http.StatusOK, recipes); err != nil {
			s.log.Error(err)
			http.Error(w, "error encoding response", http.StatusInternalServerError)
			return
		}
	}
}
