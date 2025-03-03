package handler

import "github.com/gofiber/fiber/v2"

type MainHandler struct{}

// Show landing page with search form.
func (h *MainHandler) Index(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{}, "_layout")
}
