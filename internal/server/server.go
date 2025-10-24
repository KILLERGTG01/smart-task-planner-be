package server

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/helmet/v2"
	"go.uber.org/zap"

	"github.com/KILLERGTG01/smart-task-planner-be/internal/config"
	"github.com/KILLERGTG01/smart-task-planner-be/internal/middleware"
	"github.com/KILLERGTG01/smart-task-planner-be/internal/routes"
)

func NewApp(cfg *config.Config) (*fiber.App, *middleware.AuthMiddleware) {
	authMiddleware := middleware.NewAuthMiddleware(cfg)

	app := fiber.New(fiber.Config{
		ReadTimeout:           30 * time.Second,
		WriteTimeout:          30 * time.Second,
		IdleTimeout:           120 * time.Second,
		BodyLimit:             4 * 1024 * 1024,
		ServerHeader:          "SmartTracker",
		AppName:               "SmartTracker API v1.0",
		ErrorHandler:          customErrorHandler,
		DisableStartupMessage: cfg.Env == "production",
	})

	app.Use(recover.New(recover.Config{
		EnableStackTrace: cfg.Env == "development",
	}))

	app.Use(helmet.New(helmet.Config{
		XSSProtection:      "1; mode=block",
		ContentTypeNosniff: "nosniff",
		XFrameOptions:      "DENY",
		HSTSMaxAge:         31536000,
		ReferrerPolicy:     "strict-origin-when-cross-origin",
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.AllowedOrigins,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Requested-With",
		AllowCredentials: false,
		MaxAge:           86400,
	}))

	app.Use(requestid.New())
	app.Use(middleware.ConfigMiddleware(cfg))

	if cfg.Env != "production" {
		app.Use(logger.New(logger.Config{
			Format: "[${time}] ${status} - ${method} ${path} - ${latency}\n",
		}))
	}

	routes.SetupRoutes(app, authMiddleware)

	zap.L().Info("routes registered successfully")
	return app, authMiddleware
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	zap.L().Error("request error",
		zap.Error(err),
		zap.Int("status", code),
		zap.String("method", c.Method()),
		zap.String("path", c.Path()),
		zap.String("ip", c.IP()),
	)

	return c.Status(code).JSON(fiber.Map{
		"error": message,
	})
}
