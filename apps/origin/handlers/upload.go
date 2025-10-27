package handlers

import (
	"encoding/json"
	"net/http"
	"path"

	"bitoracdn/origin/db"
	"bitoracdn/origin/storage"
)

type UploadHandler struct {
	DB      *db.DB
	Storage storage.Storage
	Bucket  string
}

type UploadInitRequest struct {
	UserID   string `json:"user_id"`
	Filename string `json:"filename"`
}

type UploadInitResponse struct {
	UploadURL string `json:"upload_url"`
	ObjectKey string `json:"object_key"`
	Bucket    string `json:"bucket"`
}

type UploadCompleteRequest struct {
	UserID      string `json:"user_id"`
	ObjectKey   string `json:"object_key"`
	ContentType string `json:"content_type"`
	Size        int64  `json:"size"`
}

// UploadInitHandler — returns Supabase upload endpoint
func (h *UploadHandler) UploadInitHandler(w http.ResponseWriter, r *http.Request) {
	var req UploadInitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	key := path.Join(req.UserID, req.Filename)
	url := h.Storage.PresignedPut(key) // only takes key, no ctx/duration

	resp := UploadInitResponse{
		UploadURL: url,
		ObjectKey: key,
		Bucket:    h.Bucket,
	}
	writeJSON(w, resp)
}

// UploadCompleteHandler — store metadata in DB
func (h *UploadHandler) UploadCompleteHandler(w http.ResponseWriter, r *http.Request) {
	var req UploadCompleteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// (Optional) you could later verify the object exists via a HEAD request to Supabase.

	asset := db.Asset{
		UserID:      req.UserID,
		Key:         req.ObjectKey,
		Bucket:      h.Bucket,
		Size:        req.Size,
		ContentType: req.ContentType,
		ETag:        "",       // not available from REST
		CacheTTL:    300,      // default 5 min cache
	}

	if err := h.DB.InsertAsset(r.Context(), asset); err != nil {
		http.Error(w, "failed to insert metadata", http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]any{
		"status": "uploaded",
		"key":    req.ObjectKey,
		"size":   req.Size,
	})
}

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
