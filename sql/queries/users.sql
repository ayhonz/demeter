-- name: CreateUser :one
INSERT INTO users (id, updated_at, created_at, first_name, last_name)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;
