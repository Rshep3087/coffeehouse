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
}

func (s *server) health() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "healthy")
	}
}

func (s *server) handleCreateRecipe() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var recip postgres.CreateRecipeParams
		if err := json.NewDecoder(r.Body).Decode(&recip); err != nil {
			http.Error(w, "invalid json body", http.StatusBadRequest)
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

		log.Info("recipe added")

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(newRecipe); err != nil {
			return
		}
	}
}

func (s *server) handleGetRecipe() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := httprouter.ParamsFromContext(r.Context())
		recipeID := params.ByName("id")

		log := s.log.With(
			"id", recipeID,
		)

		rid, err := strconv.Atoi(recipeID)
		if err != nil {
			log.Error(err)
			http.Error(w, "invalid id param", http.StatusBadRequest)
			return
		}

		log.Info("getting recipe")

		recipe, err := s.queries.GetRecipe(r.Context(), int64(rid))
		if err != nil {
			log.Error(err)
			http.Error(w, "error getting recipe", http.StatusInternalServerError)
			return
		}

		log.Info("recipe got")

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(recipe); err != nil {
			return
		}
	}
}
