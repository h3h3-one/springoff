package album

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"log/slog"
	"springoff/internal/models/album"
)

type Album struct {
	album *album.Album
}

func New(db *sql.DB) *Album {
	alb := album.New(db)
	return &Album{album: alb}
}

func (a *Album) GetAlbum() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		idAlbum := c.Params("idAlbum")
		albumImages, err := a.album.GetAlbumImages(idAlbum)
		if err != nil {
			slog.Error("Error get album images", "error", err, "ID album", idAlbum)
			return c.SendStatus(404)
		}

		title, err := a.album.GetTitleAlbum(idAlbum)
		if err != err {
			slog.Error("Error get title album", "error", err, "id album", idAlbum)
		}

		return c.Render("album",
			fiber.Map{"AlbumImages": albumImages,
				"Title": title})
	}
}

func (a *Album) UploadAlbum() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {

		return c.Redirect("//localhost:8080/")
	}
}
