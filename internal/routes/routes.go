package routes

import (
	"github.com/KILLERGTG01/smart-task-planner-be/internal/handlers"
	"github.com/KILLERGTG01/smart-task-planner-be/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, authMiddleware *middleware.AuthMiddleware) {
	SetupHealthRoutes(app)
	SetupAuthRoutes(app, authMiddleware)

	// Public API routes (no auth required)
	api := app.Group("/api")
	api.Post("/generate", handlers.GenerateHandler)
	api.Post("/generate/stream", handlers.GenerateStreamHandler)

	// Protected API routes (auth required)
	protectedAPI := app.Group("/api", authMiddleware.AuthRequired())
	protectedAPI.Get("/history", handlers.HistoryHandler)
}
