package album

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log/slog"
	"os"
	"springoff/internal/models/upload"
	"springoff/internal/util"
	"strings"
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

func (a *Album) GetAlbums() ([]Albums, error) {
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

func (a *Album) Upload(upload *upload.Upload, c *fiber.Ctx) error {
	albumName := util.RandomName()
	coverName := util.RandomName()
	pathImage := make([]string, len(upload.Album))
	slog.Info("save cover", "name", upload.Cover[0].Filename, "size", upload.Cover[0].Size)

	//creating a folder within a folder with "covers"
	if err := os.Mkdir(fmt.Sprintf("%s%s", PathCovers, albumName), 0777); err != nil {
		slog.Error("Error create cover dir", "err", err, "album name", albumName, "path", fmt.Sprintf("%s%s", PathCovers, albumName))
		return fmt.Errorf("error create cover dir")
	}
	cvrCreate, err := os.Create(fmt.Sprintf("%s%s/%s.jpg", PathCovers, albumName, coverName))
	if err != nil {
		return fmt.Errorf("error create cover image")
	}
	if err := util.WriteFile(upload.Cover[0], cvrCreate); err != nil {
		return err
	}
	err = cvrCreate.Close()
	if err != nil {
		return err
	}
	//creating a folder within a folder with "albums"
	if err := os.Mkdir(fmt.Sprintf("%s%s", PathAlbums, albumName), 0777); err != nil {
		slog.Error("Error create album dir", "err", err, "album name", albumName, "path", fmt.Sprintf("%s%s", PathAlbums, albumName))
		return fmt.Errorf("error create album dir")
	}
	//adding a file to the "albums" folder
	for i := 0; i < len(upload.Album); i++ {
		albumImageName := util.RandomName()
		albCreate, err := os.Create(fmt.Sprintf("%s%s/%s.jpg", PathAlbums, albumName, albumImageName))
		if err != nil {
			return fmt.Errorf("error create cover image")
		}
		if err := util.WriteFile(upload.Album[i], albCreate); err != nil {
			return err
		}
		pathImage[i] = fmt.Sprintf("%s/%s.jpg", albumName, albumImageName)
		err = albCreate.Close()
		if err != nil {
			return err
		}
	}

	slog.Info("transaction initialization")
	tx, err := a.db.Begin()

	alb, err := tx.Exec("INSERT INTO albums(title_album, path_cover) VALUES (?,?)", upload.Title[0], fmt.Sprintf("%s/%s.jpg", albumName, coverName))
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

	album := Albums{}
	row := a.db.QueryRow("SELECT path_cover FROM albums WHERE id_album=?", id)
	err := row.Scan(&album.PathCover)
	if err != nil {
		slog.Error("error query \"SELECT path_cover FROM albums WHERE id_album=?\"")
		return err
	}

	albumName := strings.Split(album.PathCover, "/")[0]
	if err := os.RemoveAll("./static/covers/" + albumName); err != nil {
		slog.Error("error delete covers from file system", "err", err)
		return err
	}
	if err := os.RemoveAll("./static/albums/" + albumName); err != nil {
		slog.Error("error delete covers from file system", "err", err)
		return err
	}

	slog.Info("album successfully deleted from file system", "album name", albumName)

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
