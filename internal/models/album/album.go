package album

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"springoff/internal/config"
	"springoff/internal/models/upload"
	"springoff/internal/util"
)

const (
	PathCovers = "./static/covers/"
	PathAlbums = "./static/albums/"
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

func (a *Album) GetAll() ([]Albums, error) {
	rows, err := a.db.Query("SELECT * FROM albums ORDER BY id_album DESC")
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

func (a *Album) GetImages(id string) ([]Images, error) {
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

func (a *Album) GetTitle(id string) (string, error) {
	row := a.db.QueryRow("SELECT title_album FROM albums WHERE id_album=?", id)

	alb := Albums{}
	err := row.Scan(&alb.TitleAlbum)
	if err != nil {
		return "", err
	}
	return alb.TitleAlbum, nil
}

func (a *Album) Upload(upload *upload.Upload, config *config.Config) error {
	//adding cover
	coverResp, err := util.UploadFile(upload.Cover[0], config)
	if err != nil {
		return err
	}
	if coverResp.Status != 200 {
		slog.Error("Cover not added", "code", coverResp.Status, "body", coverResp.Data, "success", coverResp.Success)
		return err
	}

	//adding album image
	pathImage := make([]string, len(upload.Album))
	for i := 0; i < len(upload.Album); i++ {
		albumResp, err := util.UploadFile(upload.Album[i], config)
		if err != nil {
			return err
		}
		if albumResp.Status != 200 {
			slog.Error("Album image not added", "code", albumResp.Status, "body", albumResp.Data, "success", albumResp.Success)
			return err
		}
		pathImage[i] = albumResp.Data.Link
	}

	slog.Info("transaction initialization")
	tx, err := a.db.Begin()

	alb, err := tx.Exec("INSERT INTO albums(title_album, path_cover) VALUES (?,?)", upload.Title[0], coverResp.Data.Link)
	if err != nil {
		tx.Rollback()
		return err
	}

	lastId, _ := alb.LastInsertId()
	for i := 0; i < len(pathImage); i++ {
		_, err = tx.Exec("INSERT INTO images(path_image, id_album) VALUES (?,?)", pathImage[i], lastId)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	slog.Info("the transaction was successful")

	return nil
}

func (a *Album) Delete(id int) error {

	slog.Info("delete album", "id album", id)

	res, err := a.db.Exec(`
	PRAGMA foreign_keys = ON;
	DELETE FROM albums WHERE id_album=?
	`, id)
	affect, err := res.RowsAffected()
	if affect == 0 {
		slog.Error("album deletion error", "id", id)
		return fmt.Errorf("the album hasn't been deleted")
	}
	if err != nil {
		slog.Error("error when executing a request to delete an album")
		return err
	}

	slog.Info("album successfully deleted from database")

	return nil
}
