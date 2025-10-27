package storage

import (
	"context"
	// "net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3Storage struct {
	Client *minio.Client
	Bucket string
}

// ObjectInfo represents object metadata.
type ObjectInfo struct {
	Size int64
	ETag string
}

// NewS3Storage initializes the MinIO (S3-compatible) client.
func NewS3Storage(endpoint, accessKey, secretKey, bucket string, useSSL bool) (*S3Storage, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	// ensure bucket exists
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return nil, err
	}
	if !exists {
		if err := client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, err
		}
	}

	return &S3Storage{
		Client: client,
		Bucket: bucket,
	}, nil
}

// Presigned PUT URL (for client upload)
func (s *S3Storage) PresignedPut(ctx context.Context, key string, expiry time.Duration) (string, error) {
	// reqParams := make(url.Values)
	u, err := s.Client.PresignedPutObject(ctx, s.Bucket, key, expiry)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

// Presigned GET URL (for edge download)
func (s *S3Storage) PresignedGet(ctx context.Context, key string, expiry time.Duration) (string, error) {
	u, err := s.Client.PresignedGetObject(ctx, s.Bucket, key, expiry, nil)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

// Check object metadata
func (s *S3Storage) StatObject(ctx context.Context, key string) (*ObjectInfo, error) {
	info, err := s.Client.StatObject(ctx, s.Bucket, key, minio.StatObjectOptions{})
	if err != nil {
		return nil, err
	}
	return &ObjectInfo{Size: info.Size, ETag: info.ETag}, nil
}
