package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type SupabaseStorage struct {
	Endpoint   string
	Bucket     string
	AccessKey  string // service role key
	HTTPClient *http.Client
}

func NewSupabaseStorage(endpoint, bucket, accessKey string) *SupabaseStorage {
	return &SupabaseStorage{
		Endpoint:   endpoint,
		Bucket:     bucket,
		AccessKey:  accessKey,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// PresignedPut — Supabase REST doesn’t use presigned uploads;
// instead, client uploads directly with Bearer token.
func (s *SupabaseStorage) PresignedPut(path string) string {
	return fmt.Sprintf("https://%s/storage/v1/object/%s/%s", s.Endpoint, s.Bucket, path)
}

// UploadFile — perform direct upload (server-side, optional)
func (s *SupabaseStorage) UploadFile(path string, data []byte, contentType string) error {
	url := fmt.Sprintf("https://%s/storage/v1/object/%s/%s", s.Endpoint, s.Bucket, path)
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.AccessKey)
	req.Header.Set("Content-Type", contentType)

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload failed: %s", string(body))
	}
	return nil
}

// PresignedGet — generate signed download link via Supabase API
func (s *SupabaseStorage) PresignedGet(path string, expiresSec int) (string, error) {
	url := fmt.Sprintf("https://%s/storage/v1/object/sign/%s/%s", s.Endpoint, s.Bucket, path)
	payload := map[string]int{"expiresIn": expiresSec}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+s.AccessKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("signing failed: %s", string(b))
	}

	var result struct {
		SignedURL string `json:"signedURL"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return fmt.Sprintf("https://%s%s", s.Endpoint, result.SignedURL), nil
}
