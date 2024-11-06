// Handlers per fiber riguardanti la parte frontend
package frontend

import (
	log "doable-go/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

// Handler for /
func Index() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		log.Info("[Frontend] index.html requested")
		return c.Render("index", fiber.Map{
			"title": "Doable web frontend",
		})
	}
}
