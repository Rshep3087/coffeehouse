package redis_test

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	redisClient "github.com/redis/go-redis/v9"
	"github.com/rshep3087/coffeehouse/cache"
	"github.com/rshep3087/coffeehouse/cache/redis"
	"github.com/rshep3087/coffeehouse/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRecipeCache(t *testing.T) {
	ctx := context.Background()
	t.Run("get and set recipe", func(t *testing.T) {
		s := miniredis.RunT(t)
		defer s.Close()

		c := redisClient.NewClient(&redisClient.Options{
			Addr: s.Addr(),
		})

		cache := redis.New(c)

		// Set recipe
		err := cache.SetRecipe(ctx, 1, &postgres.Recipe{
			ID:         1,
			RecipeName: "test recipe",
			BrewMethod: postgres.BrewMethodChemex,
			WeightUnit: postgres.WeightUnitG,
		})
		require.NoError(t, err)

		want := &postgres.Recipe{
			ID:         1,
			RecipeName: "test recipe",
			BrewMethod: postgres.BrewMethodChemex,
			WeightUnit: postgres.WeightUnitG,
		}

		// Get recipe
		got, err := cache.GetRecipe(ctx, 1)
		require.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("get recipe miss", func(t *testing.T) {
		s := miniredis.RunT(t)
		defer s.Close()

		c := redisClient.NewClient(&redisClient.Options{
			Addr: s.Addr(),
		})

		r := redis.New(c)

		_, err := r.GetRecipe(ctx, 1)
		assert.ErrorIs(t, err, cache.ErrCacheMiss)
	})
}
