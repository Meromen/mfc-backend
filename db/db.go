package db

import (
	"database/sql"
	_ "github.com/jackc/pgx"
	_ "github.com/lib/pq"
)

var defaultPgUrl = "postgres://postgres@127.0.0.1:5432/mfc?sslmode=disable"

func Connect(connStr *string) (*sql.DB, error) {
	if connStr == nil {
		connStr = &defaultPgUrl
	}

	db, err := sql.Open("postgres", *connStr)
	return db, err
}

type Storage interface {
	SelectAll() ([]DBRow, error)
	UpdateAll([]DBRow) error
}

type DBRow interface {
	GetId () string
}