package db

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

var Pool *pgxpool.Pool

func Connect(databaseURL string) {
	cfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		zap.L().Fatal("parse db config", zap.Error(err))
	}
	cfg.MaxConns = 20
	cfg.MinConns = 1
	cfg.MaxConnIdleTime = 10 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		zap.L().Fatal("db connection failed", zap.Error(err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := pool.Ping(ctx); err != nil {
		zap.L().Fatal("db ping failed", zap.Error(err))
	}

	Pool = pool
	zap.L().Info("database connected successfully")
}

func Migrate(databaseURL string) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		zap.L().Fatal("parse db config", zap.Error(err))
	}
	db := stdlib.OpenDB(*config.ConnConfig)
	defer db.Close()

	driver, err := pgx.WithInstance(db, &pgx.Config{})
	if err != nil {
		zap.L().Fatal("migration driver error", zap.Error(err))
	}

	migrationPath, err := filepath.Abs("migrations")
	if err != nil {
		zap.L().Fatal("migration path error", zap.Error(err))
	}

	if _, err := os.Stat(migrationPath); os.IsNotExist(err) {
		zap.L().Info("migrations directory not found, skipping migrations", zap.String("path", migrationPath))
		return
	}

	fileSource, err := (&file.File{}).Open("file://" + migrationPath)
	if err != nil {
		zap.L().Fatal("migration source error", zap.Error(err))
	}

	m, err := migrate.NewWithInstance("file", fileSource, "pgx", driver)
	if err != nil {
		zap.L().Fatal("migration init error", zap.Error(err))
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		zap.L().Fatal("migration up error", zap.Error(err))
	} else {
		zap.L().Info("migrations applied successfully")
	}
}

func Close() {
	if Pool != nil {
		Pool.Close()
	}
}
