package main

import (
	"github.com/KILLERGTG01/smart-task-planner-be/internal/config"
	"github.com/KILLERGTG01/smart-task-planner-be/internal/db"
	"github.com/KILLERGTG01/smart-task-planner-be/internal/logger"
)

func main() {
	cfg := config.Load()

	if err := logger.Init(cfg.Env); err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
	defer logger.Sync()

	db.Migrate(cfg.DatabaseURL)
}
