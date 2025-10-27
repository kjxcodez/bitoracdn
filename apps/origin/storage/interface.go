package storage


// Storage defines the methods common to all storage backends.
type Storage interface {
	PresignedPut(path string) string
	PresignedGet(path string, expiresSec int) (string, error)
	UploadFile(path string, data []byte, contentType string) error
}
