package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"path"
	"time"

	"bitoracdn/origin/db"
	"bitoracdn/origin/storage"
)

type UploadHandler struct {
	DB      *db.DB
	Storage *storage.S3Storage
	Bucket  string
}

// UploadInitRequest is sent by the Control Plane
type UploadInitRequest struct {
	UserID   string `json:"user_id"`
	Filename string `json:"filename"`
}

// UploadInitResponse is returned by Origin
type UploadInitResponse struct {
	PresignedURL string `json:"presigned_url"`
	ObjectKey    string `json:"object_key"`
	Bucket       string `json:"bucket"`
}

// UploadCompleteRequest confirms upload success
type UploadCompleteRequest struct {
	UserID      string `json:"user_id"`
	ObjectKey   string `json:"object_key"`
	ContentType string `json:"content_type"`
}