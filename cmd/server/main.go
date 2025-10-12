package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/KILLERGTG01/smart-task-planner-be/internal/config"
	"github.com/KILLERGTG01/smart-task-planner-be/internal/db"
	"github.com/KILLERGTG01/smart-task-planner-be/internal/logger"
	"github.com/KILLERGTG01/smart-task-planner-be/internal/server"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()

	// Initialize logger
	if err := logger.Init(cfg.Env); err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
	defer logger.Sync()

	// connect DB and run migrations
	db.ConnectAndMigrate(cfg.DatabaseURL)

	app, authMiddleware := server.NewApp(cfg)

	// run server
	go func() {
		addr := ":" + cfg.Port
		zap.L().Info("starting server", zap.String("address", addr))
		if err := app.Listen(addr); err != nil {
			zap.L().Fatal("server listen failed", zap.Error(err))
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.L().Info("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Cleanup auth middleware
	authMiddleware.Cleanup()

	if err := app.ShutdownWithContext(ctx); err != nil {
		zap.L().Error("error during shutdown", zap.Error(err))
	}
	db.Close()
	zap.L().Info("server stopped")
}
