package database

import "github.com/jmoiron/sqlx"

type Storage struct {
	DB *sqlx.DB
}

func NewStorage(db *sqlx.DB) *Storage {
	return &Storage{DB: db}
}
