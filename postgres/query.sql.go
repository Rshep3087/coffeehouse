// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package postgres

import (
	"context"
	"database/sql"
	"time"
)

const createRecipe = `-- name: CreateRecipe :one
INSERT INTO public.recipes (
    recipe_name, brew_method, coffee_weight, weight_unit, grind_size, water_weight, water_unit, water_temp, water_temp_unit
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING id, recipe_name, brew_method, coffee_weight, weight_unit, grind_size, water_weight, water_unit, water_temp, water_temp_unit
`

type CreateRecipeParams struct {
	RecipeName    string          `json:"recipe_name"`
	BrewMethod    BrewMethod      `json:"brew_method"`
	CoffeeWeight  float64         `json:"coffee_weight"`
	WeightUnit    WeightUnit      `json:"weight_unit"`
	GrindSize     int32           `json:"grind_size"`
	WaterWeight   float64         `json:"water_weight"`
	WaterUnit     string          `json:"water_unit"`
	WaterTemp     sql.NullFloat64 `json:"water_temp"`
	WaterTempUnit sql.NullString  `json:"water_temp_unit"`
}

func (q *Queries) CreateRecipe(ctx context.Context, arg CreateRecipeParams) (Recipe, error) {
	row := q.db.QueryRowContext(ctx, createRecipe,
		arg.RecipeName,
		arg.BrewMethod,
		arg.CoffeeWeight,
		arg.WeightUnit,
		arg.GrindSize,
		arg.WaterWeight,
		arg.WaterUnit,
		arg.WaterTemp,
		arg.WaterTempUnit,
	)
	var i Recipe
	err := row.Scan(
		&i.ID,
		&i.RecipeName,
		&i.BrewMethod,
		&i.CoffeeWeight,
		&i.WeightUnit,
		&i.GrindSize,
		&i.WaterWeight,
		&i.WaterUnit,
		&i.WaterTemp,
		&i.WaterTempUnit,
	)
	return i, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO public.users (
    name, email, password_hash
) VALUES (
  $1, $2, $3
) RETURNING id, created_at, version
`

type CreateUserParams struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash []byte `json:"password_hash"`
}

type CreateUserRow struct {
	ID        int32     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Version   int32     `json:"version"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (CreateUserRow, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Name, arg.Email, arg.PasswordHash)
	var i CreateUserRow
	err := row.Scan(&i.ID, &i.CreatedAt, &i.Version)
	return i, err
}

const getRecipe = `-- name: GetRecipe :one
SELECT id, recipe_name, brew_method, coffee_weight, weight_unit, grind_size, water_weight, water_unit, water_temp, water_temp_unit FROM public.recipes
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetRecipe(ctx context.Context, id int64) (Recipe, error) {
	row := q.db.QueryRowContext(ctx, getRecipe, id)
	var i Recipe
	err := row.Scan(
		&i.ID,
		&i.RecipeName,
		&i.BrewMethod,
		&i.CoffeeWeight,
		&i.WeightUnit,
		&i.GrindSize,
		&i.WaterWeight,
		&i.WaterUnit,
		&i.WaterTemp,
		&i.WaterTempUnit,
	)
	return i, err
}

const getUserById = `-- name: GetUserById :one
SELECT id, created_at, name, email, password_hash, active, version FROM public.users
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUserById(ctx context.Context, id int32) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserById, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Name,
		&i.Email,
		&i.PasswordHash,
		&i.Active,
		&i.Version,
	)
	return i, err
}

const getUserRecipes = `-- name: GetUserRecipes :many
SELECT id, recipe_name, brew_method, coffee_weight, weight_unit, grind_size, water_weight, water_unit, water_temp, water_temp_unit, user_id, recipe_id, created_at FROM public.recipes
JOIN public.saved_recipes ON public.recipes.id = public.saved_recipes.recipe_id
WHERE public.saved_recipes.user_id = $1
`

type GetUserRecipesRow struct {
	ID            int64           `json:"id"`
	RecipeName    string          `json:"recipe_name"`
	BrewMethod    BrewMethod      `json:"brew_method"`
	CoffeeWeight  float64         `json:"coffee_weight"`
	WeightUnit    WeightUnit      `json:"weight_unit"`
	GrindSize     int32           `json:"grind_size"`
	WaterWeight   float64         `json:"water_weight"`
	WaterUnit     string          `json:"water_unit"`
	WaterTemp     sql.NullFloat64 `json:"water_temp"`
	WaterTempUnit sql.NullString  `json:"water_temp_unit"`
	UserID        int32           `json:"user_id"`
	RecipeID      int32           `json:"recipe_id"`
	CreatedAt     time.Time       `json:"created_at"`
}

func (q *Queries) GetUserRecipes(ctx context.Context, userID int32) ([]GetUserRecipesRow, error) {
	rows, err := q.db.QueryContext(ctx, getUserRecipes, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUserRecipesRow
	for rows.Next() {
		var i GetUserRecipesRow
		if err := rows.Scan(
			&i.ID,
			&i.RecipeName,
			&i.BrewMethod,
			&i.CoffeeWeight,
			&i.WeightUnit,
			&i.GrindSize,
			&i.WaterWeight,
			&i.WaterUnit,
			&i.WaterTemp,
			&i.WaterTempUnit,
			&i.UserID,
			&i.RecipeID,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listRecipes = `-- name: ListRecipes :many
SELECT id, recipe_name, brew_method, coffee_weight, weight_unit, grind_size, water_weight, water_unit, water_temp, water_temp_unit FROM public.recipes
ORDER BY brew_method
`

func (q *Queries) ListRecipes(ctx context.Context) ([]Recipe, error) {
	rows, err := q.db.QueryContext(ctx, listRecipes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Recipe
	for rows.Next() {
		var i Recipe
		if err := rows.Scan(
			&i.ID,
			&i.RecipeName,
			&i.BrewMethod,
			&i.CoffeeWeight,
			&i.WeightUnit,
			&i.GrindSize,
			&i.WaterWeight,
			&i.WaterUnit,
			&i.WaterTemp,
			&i.WaterTempUnit,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const saveRecipe = `-- name: SaveRecipe :exec
INSERT INTO public.saved_recipes (
    user_id, recipe_id
) VALUES (
  $1, $2
)
`

type SaveRecipeParams struct {
	UserID   int32 `json:"user_id"`
	RecipeID int32 `json:"recipe_id"`
}

func (q *Queries) SaveRecipe(ctx context.Context, arg SaveRecipeParams) error {
	_, err := q.db.ExecContext(ctx, saveRecipe, arg.UserID, arg.RecipeID)
	return err
}
