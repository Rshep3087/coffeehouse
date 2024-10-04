package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/rshep3087/coffeehouse/postgres"
	"golang.org/x/crypto/bcrypt"
)

func (s *server) routes() {
	// Health check route
	s.router.HandlerFunc(http.MethodGet, "/health", s.health())

	// Recipe routes
	s.router.HandlerFunc(http.MethodPost, "/v1/recipes", loggingmw(s.log, s.handleCreateRecipe()))
	s.router.HandlerFunc(http.MethodGet, "/v1/recipes/:id", loggingmw(s.log, s.handleGetRecipe()))
	s.router.HandlerFunc(http.MethodGet, "/v1/recipes", loggingmw(s.log, s.handleGetRecipes()))

	// User routes
	s.router.HandlerFunc(http.MethodPost, "/v1/users", loggingmw(s.log, s.handleCreateUser()))
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

func hashPassword(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}

func (s *server) handleCreateUser() http.HandlerFunc {
	// req is a struct that represents the JSON body of a POST request to create a user
	type req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password []byte `json:"password"`
	}

	// resp is a struct that represents the JSON response to a POST request to create a user
	type resp struct {
		Name      string    `json:"name"`
		Email     string    `json:"email"`
		ID        int32     `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		Version   int32     `json:"version"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		user, err := decode[req](r)
		if err != nil {
			http.Error(w, "error decoding user", http.StatusBadRequest)
			return
		}

		log := s.log.With("name", user.Name, "email", user.Email)

		log.Info("creating user")

		ha, err := hashPassword(user.Password)
		if err != nil {
			log.Error(err)
			http.Error(w, "error hashing password", http.StatusInternalServerError)
			return
		}

		newUser, err := s.queries.CreateUser(r.Context(), postgres.CreateUserParams{
			Name:         user.Name,
			Email:        user.Email,
			PasswordHash: ha,
		})
		if err != nil {
			log.Error(err)
			http.Error(w, "error creating user", http.StatusInternalServerError)
			return
		}

		log.Info("user created")

		if err = encode(w, http.StatusCreated, resp{
			Name:      user.Name,
			Email:     user.Email,
			ID:        newUser.ID,
			CreatedAt: newUser.CreatedAt,
			Version:   newUser.Version,
		}); err != nil {
			log.Error(err)
			http.Error(w, "error encoding response", http.StatusInternalServerError)
			return
		}
	}
}
