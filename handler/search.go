package handler

import (
	"github.com/Innoractive/chippiphone/entity"
	"github.com/Innoractive/chippiphone/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type SearchHandler struct{}

// Perform Search and output partial HTML (display result in table form).
// This output is used by HTMX to embed in parent page.
func (h *SearchHandler) Index(c *fiber.Ctx) error {
	// TODO implement search logic
	searchService := service.SearchService{}
	businesses, err := searchService.Search(c.Query("query"))
	if err != nil {
		log.Warnf("Unable to perform search: %v", err)
		return c.Render("result", fiber.Map{
			"businesses": []entity.Business{},
		})
	}

	return c.Render("result", fiber.Map{
		"Businesses": businesses,
	})
}
