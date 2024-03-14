package sqlite3

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

func New(path string) (*sql.DB, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("path not exist: %w", err)
	}

	//Create connection for database
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("fail to connection db: %w", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(`
		PRAGMA foreign_keys=on;

		CREATE TABLE IF NOT EXISTS albums(
			id_album INTEGER PRIMARY KEY AUTOINCREMENT,
			title_album VARCHAR(30) NOT NULL,
			path_cover VARCHAR(30) NOT NULL
		);

		CREATE TABLE IF NOT EXISTS images(
			id_image INTEGER PRIMARY KEY AUTOINCREMENT,
			path_image VARCHAR(30) NOT NULL,
			id_album VARCHAR(30) NOT NULL,
			CONSTRAINT fk_album
			FOREIGN KEY (id_album)
			REFERENCES albums (id_album)
			ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS users(
			id_user INTEGER PRIMARY KEY AUTOINCREMENT,
			login VARCHAR(30) NOT NULL,
			password VARCHAR(30) NOT NULL
		);

		CREATE TABLE IF NOT EXISTS uuid(
			id_uuid INTEGER PRIMARY KEY AUTOINCREMENT,
			uuid VARCHAR(30) NOT NULL
		);

		CREATE UNIQUE INDEX IF NOT EXISTS index_uuid
		on uuid (uuid);
	`)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return db, nil
}
