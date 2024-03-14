package contact

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"springoff/internal/models/contact"
)

type Contact struct {
	contact *contact.Contact
}

func New(db *sql.DB) *Contact {
	cont := contact.New(db)
	return &Contact{contact: cont}
}

func (c *Contact) GetContact() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.Render("contact", fiber.Map{"Title": "Контакты", "Contact": "active"})
	}
}
