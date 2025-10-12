package routes

import (
	"github.com/KILLERGTG01/smart-task-planner-be/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, authMiddleware *middleware.AuthMiddleware) {
	SetupHealthRoutes(app)
	api := app.Group("/api", authMiddleware.AuthRequired())
	SetupPlanRoutes(api)
}
