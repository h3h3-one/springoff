package storage

import "database/sql"

type Storage struct {
	db *sql.DB
}

func New(path string) (*sql.DB, err) {
	
}
