package service

import "database/sql"

type Service struct {
	db *sql.DB
}

func New(storage *sql.DB) *Service {
	return &Service{db: storage}
}
