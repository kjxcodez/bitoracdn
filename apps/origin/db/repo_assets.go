package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Asset represents a record in origin_meta.assets table.
type Asset struct {
	ID          uuid.UUID
	UserID      string
	Key         string
	Bucket      string
	Size        int64
	ContentType string
	ETag        string
	CacheTTL    int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// InsertAsset adds a new file metadata record into the DB.
func (db *DB) InsertAsset(ctx context.Context, a Asset) error {
	query := `
		INSERT INTO origin_meta.assets (user_id, key, bucket, size, content_type, etag, cache_ttl)
		VALUES ($1, $2, $3, $4, $5, $6, COALESCE($7, 300))
	`
	_, err := db.Pool.Exec(ctx, query,
		a.UserID,
		a.Key,
		a.Bucket,
		a.Size,
		a.ContentType,
		a.ETag,
		a.CacheTTL,
	)
	return err
}

// GetAssetByKey fetches a file record using its key.
func (db *DB) GetAssetByKey(ctx context.Context, key string) (*Asset, error) {
	query := `
		SELECT id, user_id, key, bucket, size, content_type, etag, cache_ttl, created_at, updated_at
		FROM origin_meta.assets
		WHERE key = $1
		LIMIT 1
	`
	row := db.Pool.QueryRow(ctx, query, key)

	var a Asset
	err := row.Scan(
		&a.ID,
		&a.UserID,
		&a.Key,
		&a.Bucket,
		&a.Size,
		&a.ContentType,
		&a.ETag,
		&a.CacheTTL,
		&a.CreatedAt,
		&a.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &a, nil
}
