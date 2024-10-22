package web

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/peterldowns/pgtestdb"
	"github.com/peterldowns/pgtestdb/migrators/golangmigrator"
	"github.com/rshep3087/coffeehouse/cache"
	"github.com/rshep3087/coffeehouse/postgres"
	chsql "github.com/rshep3087/coffeehouse/sql"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func NewDB(t *testing.T) *sql.DB {
	t.Helper()

	gm := golangmigrator.New(
		"schema",
		golangmigrator.WithFS(chsql.SchemaFS),
	)

	conf := pgtestdb.Config{
		DriverName: "postgres",
		Host:       "localhost",
		Port:       "5432",
		User:       "user",
		Password:   "password",
		Database:   "coffeehousedb",
		Options:    "sslmode=disable",
	}

	var migrator pgtestdb.Migrator = gm
	return pgtestdb.New(t, conf, migrator)
}

func TestHandleCreateRecipe(t *testing.T) {
	ctx := context.Background()
	log := zap.NewExample(zap.Development()).Sugar()

	t.Run("ok", func(t *testing.T) {
		t.Parallel()
		db := NewDB(t)
		defer db.Close()

		psMock := &PubSubMock{
			PublishFunc: func(topic string, data []byte) error {
				return nil
			},
		}
		cacheMock := &cache.RecipeCacherMock{}

		s := NewServer(log, psMock, cacheMock)
		s.Queries = postgres.New(db)

		// create a user
		_, err := s.Queries.CreateUser(ctx, postgres.CreateUserParams{
			Name:         "test user",
			Email:        "user@email.com",
			PasswordHash: []byte("password"),
		})
		require.NoError(t, err)

		// create a recipe
		w := httptest.NewRecorder()

		payload, err := os.ReadFile("testdata/create_recipe.json")
		require.NoError(t, err)

		r := httptest.NewRequest(http.MethodPost, "/v1/recipes", bytes.NewReader(payload))
		s.ServeHTTP(w, r)

		require.Equal(t, http.StatusCreated, w.Code)

		var got postgres.Recipe
		require.NoError(t, json.NewDecoder(w.Body).Decode(&got))
		want := postgres.Recipe{
			ID:           1,
			RecipeName:   "sample",
			BrewMethod:   postgres.BrewMethodChemex,
			WeightUnit:   postgres.WeightUnitG,
			GrindSize:    21,
			WaterWeight:  500,
			WaterUnit:    "g",
			CoffeeWeight: 20,
			UserID:       1,
		}
		require.Equal(t, want, got)
	})
}

func TestHandleGetRecipe(t *testing.T) {
	ctx := context.Background()
	log := zap.NewExample(zap.Development()).Sugar()

	t.Run("ok", func(t *testing.T) {
		t.Parallel()
		db := NewDB(t)
		defer db.Close()

		psMock := &PubSubMock{}
		cacheMock := &cache.RecipeCacherMock{
			GetRecipeFunc: func(ctx context.Context, id int64) (*postgres.Recipe, error) {
				return nil, cache.ErrCacheMiss
			},
			SetRecipeFunc: func(ctx context.Context, id int64, recipe *postgres.Recipe) error {
				return nil
			},
		}

		s := NewServer(log, psMock, cacheMock)
		s.Queries = postgres.New(db)

		_, err := s.Queries.CreateUser(ctx, postgres.CreateUserParams{
			Name:         "test user",
			Email:        "user@email.com",
			PasswordHash: []byte("password"),
		})
		require.NoError(t, err)

		// create a recipe
		_, err = s.Queries.CreateRecipe(ctx, postgres.CreateRecipeParams{
			RecipeName: "test recipe",
			BrewMethod: postgres.BrewMethodChemex,
			WeightUnit: postgres.WeightUnitG,
			UserID:     1,
		})
		require.NoError(t, err)

		// get the recipe
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/v1/recipes/1", nil)
		s.ServeHTTP(w, r)

		require.Equal(t, http.StatusOK, w.Code)
		var got postgres.Recipe
		require.NoError(t, json.NewDecoder(w.Body).Decode(&got))
		want := postgres.Recipe{
			ID:          1,
			RecipeName:  "test recipe",
			BrewMethod:  postgres.BrewMethodChemex,
			WeightUnit:  postgres.WeightUnitG,
			GrindSize:   0,
			WaterWeight: 0,
			WaterUnit:   "",
			UserID:      1,
		}
		require.Equal(t, want, got)
	})
}
