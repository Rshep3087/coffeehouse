package cache

import (
	"context"
	"errors"

	"github.com/rshep3087/coffeehouse/postgres"
)

var ErrCacheMiss = errors.New("cache miss")

// RecipeCacher is an interface for caching recipes
//
//go:generate moq -rm -out recipe_cacher_moq.go . RecipeCacher
type RecipeCacher interface {
	// GetRecipe returns a recipe from the cache if it exists
	// If it does not exist, it returns ErrCacheMiss
	GetRecipe(ctx context.Context, id int64) (*postgres.Recipe, error)

	// SetRecipe sets a recipe in the cache
	SetRecipe(ctx context.Context, id int64, recipe *postgres.Recipe) error
}
