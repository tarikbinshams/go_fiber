package utils

import (
	"fmt"
	"mime/multipart"
	"strings"
)

// ValidateFile checks file type and size
func ValidateFile(file *multipart.FileHeader, acceptedTypes []string, maxSize int64) error {
	// Check if file exists
	if file == nil {
		return fmt.Errorf("no file uploaded")
	}

	// Check file size
	if file.Size > maxSize {
		return fmt.Errorf("file size exceeds %d bytes", maxSize)
	}

	// Extract file extension
	fileExt := strings.ToLower(file.Filename[strings.LastIndex(file.Filename, ".")+1:])

	// Validate file type
	isValidType := false
	for _, t := range acceptedTypes {
		if fileExt == t {
			isValidType = true
			break
		}
	}
	if !isValidType {
		return fmt.Errorf("invalid file type, allowed: %s", strings.Join(acceptedTypes, ", "))
	}

	return nil
}
