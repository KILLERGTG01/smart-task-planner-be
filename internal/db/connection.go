package db

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

var Pool *pgxpool.Pool

func ConnectAndMigrate(databaseURL string) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		log.Fatalf("parse db config: %v", err)
	}
	db := stdlib.OpenDB(*config.ConnConfig)
	defer db.Close()

	driver, err := pgx.WithInstance(db, &pgx.Config{})
	if err != nil {
		log.Fatalf("migration driver error: %v", err)
	}

	migrationPath, err := filepath.Abs("migrations")
	if err != nil {
		log.Fatalf("migration path error: %v", err)
	}

	if _, err := os.Stat(migrationPath); os.IsNotExist(err) {
		log.Printf("migrations directory not found at %s, skipping migrations", migrationPath)
	} else {
		// Create file source
		fileSource, err := (&file.File{}).Open("file://" + migrationPath)
		if err != nil {
			log.Fatalf("migration source error: %v", err)
		}

		// Create migrate instance
		m, err := migrate.NewWithInstance("file", fileSource, "pgx", driver)
		if err != nil {
			log.Fatalf("migration init error: %v", err)
		}

		// Run migrations
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("migration up error: %v", err)
		} else {
			log.Println("migrations applied successfully")
		}
	}

	// Now create the connection pool
	cfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		log.Fatalf("parse db config: %v", err)
	}
	cfg.MaxConns = 20
	cfg.MinConns = 1
	cfg.MaxConnIdleTime = 10 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		log.Fatalf("db connection failed: %v", err)
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("db ping failed: %v", err)
	}

	Pool = pool
	log.Println("database connected successfully")
}

func Close() {
	if Pool != nil {
		Pool.Close()
	}
}
