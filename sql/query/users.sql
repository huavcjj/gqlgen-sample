-- name: ListUsersByIDs :many
SELECT id, name
FROM users
WHERE id = ANY ($1::text[]);

-- name: GetUser :one
SELECT id, name
FROM users
WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (id, name)
VALUES ($1, $2)
RETURNING id, name;