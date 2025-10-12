package routes

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func SetupHealthRoutes(app *fiber.App) {
	app.Get("/health", healthCheckHandler)
}

func healthCheckHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":    "ok",
		"timestamp": time.Now().Unix(),
		"service":   "smart-task-planner-api",
	})
}
