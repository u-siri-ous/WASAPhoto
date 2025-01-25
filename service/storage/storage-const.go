package storage

const (
	BasePath = "/tmp/storage"

	UploadedPhotoFolder = "/photos"
	AllowedMimeType     = "image/jpeg"

	MaxRequestFileSize = 12 << 20  // 12 MB
	MaxFileSize        = 10 << 20  // 10 MB
)
