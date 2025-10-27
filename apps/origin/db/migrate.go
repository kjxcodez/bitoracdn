package db

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations runs all pending DB migrations
func RunMigrations() error {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return fmt.Errorf("missing DATABASE_URL")
	}

	// Get absolute path for local migrations folder
	path, err := filepath.Abs("db/migrations")
	if err != nil {
		return err
	}

	m, err := migrate.New(
		fmt.Sprintf("file://%s", path),
		dbURL,
	)
	if err != nil {
		return fmt.Errorf("migration init failed: %w", err)
	}

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("✅ Database is up-to-date.")
			return nil
		}
		return fmt.Errorf("migration failed: %w", err)
	}

	log.Println("✅ Database migrations applied successfully.")
	return nil
}
