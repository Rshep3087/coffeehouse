-- name: GetRecipe :one
SELECT * FROM public.recipes
WHERE id = $1 LIMIT 1;

-- name: ListRecipes :many
SELECT * FROM public.recipes
ORDER BY brew_method;

-- name: CreateRecipe :one
INSERT INTO public.recipes (
    recipe_name, brew_method, coffee_weight, weight_unit, grind_size, water_weight, water_unit, water_temp, water_temp_unit
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: CreateUser :one
INSERT INTO public.users (
    name, email, password_hash
) VALUES (
  $1, $2, $3
) RETURNING id, created_at, version;
