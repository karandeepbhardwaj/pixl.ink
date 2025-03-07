package image

import (
	"fmt"
	"mime/multipart"
	"strings"
)

var allowedTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

func Validate(header *multipart.FileHeader, maxSize int64) error {
	if header.Size > maxSize {
		return fmt.Errorf("file too large: %d bytes (max %d)", header.Size, maxSize)
	}
	contentType := header.Header.Get("Content-Type")
	if !allowedTypes[contentType] {
		return fmt.Errorf("unsupported file type: %s", contentType)
	}
	return nil
}

func GetExtension(contentType string) string {
	parts := strings.Split(contentType, "/")
	if len(parts) == 2 {
		return "." + parts[1]
	}
	return ".bin"
}
