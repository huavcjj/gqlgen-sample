// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0

package db

import ()

type Todo struct {
	ID     string
	Text   string
	Done   bool
	UserID string
}

type User struct {
	ID   string
	Name string
}
