package utils

import (
	"mime"
	"mime/multipart"
	"path/filepath"
)

func GetFileContentType(file *multipart.FileHeader) string {
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = GuessContentType(file.Filename)
	}
	return contentType
}

func GuessContentType(name string) string {
	contentType := mime.TypeByExtension(filepath.Ext(name))
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	return contentType
}
