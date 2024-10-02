package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
	"github.com/rshep3087/coffeehouse/cache"
	"github.com/rshep3087/coffeehouse/postgres"
)

var _ cache.RecipeCacher = (*Cache)(nil)

type RecipeCacheKey struct {
	ID int64
}

func (r RecipeCacheKey) String() string {
	return "recipe:" + strconv.FormatInt(r.ID, 10)
}

// Cache is a cache for recipes
type Cache struct {
	client *redis.Client
}

func New(client *redis.Client) *Cache {
	return &Cache{client: client}
}

// GetRecipe returns a recipe from the cache if it exists
// If it does not exist, it returns ErrCacheMiss
func (c *Cache) GetRecipe(ctx context.Context, id int64) (*postgres.Recipe, error) {
	cacheKey := RecipeCacheKey{ID: id}
	// get the recipe from the cache
	recipe, err := c.client.Get(ctx, cacheKey.String()).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, cache.ErrCacheMiss
		}
		return nil, fmt.Errorf("getting recipe from cache: %w", err)
	}

	// unmarshal the recipe
	var r postgres.Recipe
	if err := json.Unmarshal([]byte(recipe), &r); err != nil {
		return nil, fmt.Errorf("unmarshaling recipe: %w", err)
	}

	return &r, nil
}

func (c *Cache) SetRecipe(ctx context.Context, id int64, recipe *postgres.Recipe) error {
	cacheKey := RecipeCacheKey{ID: id}
	recipeJSON, err := json.Marshal(recipe)
	if err != nil {
		return fmt.Errorf("marshaling recipe: %w", err)
	}

	// set the recipe in the cache
	if err := c.client.Set(ctx, cacheKey.String(), recipeJSON, 0).Err(); err != nil {
		return fmt.Errorf("setting recipe in cache: %w", err)
	}

	return nil
}
