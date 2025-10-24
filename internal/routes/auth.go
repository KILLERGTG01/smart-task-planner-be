package routes

import (
	"github.com/KILLERGTG01/smart-task-planner-be/internal/handlers"
	"github.com/KILLERGTG01/smart-task-planner-be/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(app *fiber.App, authMiddleware *middleware.AuthMiddleware) {
	auth := app.Group("/auth")

	auth.Get("/login", handlers.LoginHandler)
	auth.Get("/callback", handlers.CallbackHandler)
	auth.Post("/exchange", handlers.ExchangeTokenHandler)
	auth.Post("/refresh", handlers.RefreshTokenHandler)
	auth.Get("/logout", handlers.LogoutHandler)

	auth.Get("/profile", authMiddleware.AuthRequired(), handlers.UserProfileHandler)
}
