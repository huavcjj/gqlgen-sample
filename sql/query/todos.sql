-- name: GetTodos :many
SELECT id, text, done, user_id
FROM todos;

-- name: CreateTodo :one
INSERT INTO todos (id, text, done, user_id)
VALUES ($1, $2, $3, $4)
RETURNING *;