-- name: CreateRecipe :one
INSERT INTO recipes (id, updated_at, created_at, title, description, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetRecipes :many
SELECT * FROM recipes
LIMIT $1
OFFSET $2;

-- name: GetRecipeByID :one
SELECT * FROM recipes WHERE id = $1;


