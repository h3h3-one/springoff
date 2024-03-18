package controllers

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"log/slog"
	"springoff/internal/config"
	"springoff/internal/models/album"
	"springoff/internal/models/upload"
	"strconv"
	"strings"
)

type Album struct {
	album *album.Album
}

func AlbumNew(db *sql.DB) *Album {
	alb := album.New(db)
	return &Album{album: alb}
}

func (a *Album) GetAlbum() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		idAlbum := c.Params("idAlbum")
		albumImages, err := a.album.GetImages(idAlbum)
		if err != nil {
			slog.Error("Error get album images", "error", err, "ID album", idAlbum)
			return c.SendStatus(404)
		}

		title, err := a.album.GetTitle(idAlbum)
		if err != err {
			slog.Error("Error get title album", "error", err, "id album", idAlbum)
		}

		return c.Render("album",
			fiber.Map{"AlbumImages": albumImages,
				"Title": title})
	}
}

func (a *Album) UploadAlbum(config *config.Config) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if form, err := c.MultipartForm(); err == nil {
			up, err := upload.New(form.Value["title"], form.File["cover"], form.File["albumImage"])
			if err != nil {
				slog.Error("Error validate forms adding album", "err", err)
				return c.Status(400).JSON(fiber.Map{
					"error": err.Error(),
				})
			}

			if err := a.album.Upload(up, config); err != nil {
				slog.Error("Error upload album", "err", err)
				return c.Status(400).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
		}

		slog.Info("Album added successfully")
		return c.Status(200).JSON(fiber.Map{
			"successfully": "album added successfully",
		})
	}
}

func (a *Album) DeleteAlbum() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		temp := string(c.Body())                                  //id=...
		idAlbum, err := strconv.Atoi(strings.Split(temp, "=")[1]) // int id
		if err != nil {
			return err
		}
		if err := a.album.Delete(idAlbum); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(200).JSON(fiber.Map{
			"successfully": "album deleted successfully",
		})
	}
}
