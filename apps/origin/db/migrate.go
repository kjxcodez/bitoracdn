package db

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib" // register "pgx" driver
)

// RunMigrations manually runs SQL files from remote URLs.
func RunMigrations() error {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return fmt.Errorf("missing DATABASE_URL")
	}

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		return fmt.Errorf("db open: %w", err)
	}
	defer db.Close()

	// 🔹 List your raw SQL files here (order matters)
	migrations := []string{
		"https://raw.githubusercontent.com/kjxcodez/bitoracdn/refs/heads/main/apps/origin/db/migrations/001_init_schema.up.sql",
		// add more as you version your schema
	}

	ctx := context.Background()

	for _, url := range migrations {
		log.Printf("🚀 Applying migration: %s", url)

		resp, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("fetch %s: %w", url, err)
		}
		defer resp.Body.Close()

		sqlBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("read %s: %w", url, err)
		}

		sqlText := string(sqlBytes)
		if _, err := db.ExecContext(ctx, sqlText); err != nil {
			return fmt.Errorf("exec %s: %w", url, err)
		}

		log.Printf("✅ Migration applied: %s", url)
	}

	log.Println("✅ All remote migrations executed successfully.")
	return nil
}
