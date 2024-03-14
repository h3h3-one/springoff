package home

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"log/slog"
	"springoff/internal/models/album"
)

type Home struct {
	album *album.Album
}

func New(db *sql.DB) *Home {
	alb := album.New(db)
	return &Home{album: alb}
}

func (h *Home) GetAllAlbums() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		allAlbums, err := h.album.GetAlbums()
		if err != nil {
			slog.Error("Error get all albums", "error", err)
		}
		return c.Render("home", fiber.Map{"AllAlbums": allAlbums, "Title": "Главная страница", "Home": "active"})
	}
}
