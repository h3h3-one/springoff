package about

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"springoff/internal/models/about"
)

type About struct {
	about *about.About
}

func New(db *sql.DB) *About {
	abt := about.New(db)
	return &About{about: abt}
}

func (a *About) GetAbout() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.Render("about", fiber.Map{"Title": "О себе", "About": "active"})
	}
}
