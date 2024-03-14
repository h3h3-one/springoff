package album

import (
	"database/sql"
	"errors"
	"log/slog"
)

type Album struct {
	db *sql.DB
}

type Albums struct {
	IdAlbum    int
	TitleAlbum string
	PathCover  string
}

type Images struct {
	IdImage   int
	PathImage string
	IdAlbum   int
}

func New(storage *sql.DB) *Album {
	return &Album{db: storage}
}

func (a *Album) GetAlbums() ([]Albums, error) {
	rows, err := a.db.Query("SELECT * FROM albums")
	if err != nil {
		return nil, err
	}

	var albums []Albums
	for rows.Next() {
		alb := Albums{}
		err := rows.Scan(&alb.IdAlbum, &alb.TitleAlbum, &alb.PathCover)
		if err != nil {
			slog.Error("Error mapping for object Album", "error", err)
			continue
		}
		albums = append(albums, alb)
	}

	return albums, nil
}

func (a *Album) GetAlbumImages(id string) ([]Images, error) {
	rows, err := a.db.Query("SELECT * FROM images WHERE id_album=?", id)
	if err != nil {
		return nil, err
	}

	var images []Images
	for rows.Next() {
		img := Images{}
		err := rows.Scan(&img.IdImage, &img.PathImage, &img.IdAlbum)
		if err != nil {
			slog.Error("Error mapping for object Images", "error", err)
			continue
		}
		images = append(images, img)
	}

	if images == nil {
		return nil, errors.New("images not found")
	}

	return images, nil
}

func (a *Album) GetTitleAlbum(id string) (string, error) {
	row := a.db.QueryRow("SELECT title_album FROM albums WHERE id_album=?", id)

	alb := Albums{}
	err := row.Scan(&alb.TitleAlbum)
	if err != nil {
		return "", err
	}
	return alb.TitleAlbum, nil
}
