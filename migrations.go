package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func runMigrations(ctx context.Context) error {
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	log.Println("Running database migrations...")

	m, err := migrate.New(
		"file://db/migrations",
		dbURL,
	)
	if err != nil {
		return fmt.Errorf("failed to initialize migrate: %w", err)
	}
	defer m.Close()

	// Check current version
	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get version: %w", err)
	}

	log.Printf("Current schema version: %d (dirty: %v)", version, dirty)

	// Apply migrations with timeout
	done := make(chan error, 1)
	go func() {
		start := time.Now()
		err := m.Up()
		duration := time.Since(start)

		if err != nil && err != migrate.ErrNoChange {
			log.Printf("Migration failed after %v: %v", duration, err)
			done <- err
			return
		}

		if err == migrate.ErrNoChange {
			log.Println("No new migrations to apply")
		} else {
			newVersion, _, _ := m.Version()
			log.Printf("Migrations completed: v%d → v%d (took %v)", version, newVersion, duration)
		}
		done <- nil
	}()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return fmt.Errorf("migration timeout: %w", ctx.Err())
	}
}
