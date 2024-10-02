package main

import (
	"context"
	"database/sql"
	"embed"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/peterldowns/pgtestdb"
	"github.com/peterldowns/pgtestdb/migrators/golangmigrator"
	"github.com/rshep3087/coffeehouse/cache"
	"github.com/rshep3087/coffeehouse/postgres"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

//go:embed sql/schema/*.sql
var schemaFS embed.FS

func NewDB(t *testing.T) *sql.DB {
	t.Helper()

	gm := golangmigrator.New(
		"sql/schema",
		golangmigrator.WithFS(schemaFS),
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

		s := newServer(psMock, cacheMock)
		s.queries = postgres.New(db)
		s.log = log

		// create a recipe
		_, err := s.queries.CreateRecipe(ctx, postgres.CreateRecipeParams{
			RecipeName: "test recipe",
			BrewMethod: postgres.BrewMethodChemex,
			WeightUnit: postgres.WeightUnitG,
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
		}
		require.Equal(t, want, got)
	})
}
