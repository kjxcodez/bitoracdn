package db

import (
	"context"
	"os"
	"time"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

func Connect() (*DB, error) {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		return nil, ErrMissingDBURL
	}

	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}

	config.MaxConns = 10
	config.MinConns = 1
	config.MaxConnLifetime = time.Hour

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return &DB{Pool: pool}, nil
}

func (db *DB) Close() {
	db.Pool.Close()
}

var ErrMissingDBURL = fmt.Errorf("missing DATABASE_URL")
