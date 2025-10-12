package routes

import (
	"github.com/KILLERGTG01/smart-task-planner-be/internal/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupPlanRoutes(api fiber.Router) {
	api.Post("/generate", handlers.GenerateHandler)
	api.Get("/history", handlers.HistoryHandler)
}
