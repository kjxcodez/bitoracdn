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

// UploadInitHandler generates presigned upload URL
func (h *UploadHandler) UploadInitHandler(w http.ResponseWriter, r *http.Request) {
	var req UploadInitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	key := path.Join(req.UserID, req.Filename)

	ctx := context.Background()
	url, err := h.Storage.PresignedPut(ctx, key, 15*time.Minute)
	if err != nil {
		http.Error(w, "failed to create presigned URL", http.StatusInternalServerError)
		return
	}

	resp := UploadInitResponse{
		PresignedURL: url,
		ObjectKey:    key,
		Bucket:       h.Bucket,
	}
	writeJSON(w, resp)
}

// UploadCompleteHandler verifies object and inserts metadata
func (h *UploadHandler) UploadCompleteHandler(w http.ResponseWriter, r *http.Request) {
	var req UploadCompleteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	info, err := h.Storage.StatObject(ctx, req.ObjectKey)
	if err != nil {
		http.Error(w, "object not found", http.StatusNotFound)
		return
	}

	asset := db.Asset{
		UserID:      req.UserID,
		Key:         req.ObjectKey,
		Bucket:      h.Bucket,
		Size:        info.Size,
		ContentType: req.ContentType,
		ETag:        info.ETag,
	}

	if err := h.DB.InsertAsset(ctx, asset); err != nil {
		http.Error(w, "failed to insert metadata", http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]any{
		"status": "uploaded",
		"key":    req.ObjectKey,
		"size":   info.Size,
	})
}

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
