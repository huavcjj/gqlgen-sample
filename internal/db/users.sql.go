// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0
// source: users.sql

package db

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (id, name)
VALUES ($1, $2)
RETURNING id, name
`

type CreateUserParams struct {
	ID   string
	Name string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser, arg.ID, arg.Name)
	var i User
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT id, name
FROM users
WHERE id = $1
`

func (q *Queries) GetUser(ctx context.Context, id string) (User, error) {
	row := q.db.QueryRow(ctx, getUser, id)
	var i User
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const listUsersByIDs = `-- name: ListUsersByIDs :many
SELECT id, name FROM users WHERE id = ANY($1::text[])
`

func (q *Queries) ListUsersByIDs(ctx context.Context, dollar_1 []string) ([]User, error) {
	rows, err := q.db.Query(ctx, listUsersByIDs, dollar_1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
