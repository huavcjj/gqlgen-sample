-- name: GetTodos :many
SELECT id, text, done, user_id
FROM todos;

-- name: CreateTodo :one
INSERT INTO todos (id, text, done, user_id)
VALUES (?, ?, ?, ?)
RETURNING id, text, done, user_id;

-- name: GetUser :one
SELECT id, name
FROM users
WHERE id = ?;