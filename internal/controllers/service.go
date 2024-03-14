package service

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"springoff/internal/models/service"
)

type Service struct {
	service *service.Service
}

func New(db *sql.DB) *Service {
	srv := service.New(db)
	return &Service{service: srv}
}

func (s *Service) GetService() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.Render("service", fiber.Map{"Title": "Прайс", "Service": "active"})
	}
}
