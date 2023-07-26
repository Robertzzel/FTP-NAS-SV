package utils

import (
	"mime"
	"net/http"
	"os"
)

type FileDetails struct {
	Name  string
	Size  int64
	IsDir bool
	Type  string
}

func GetFileType(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", nil
	}
	defer file.Close()

	mimeType := mime.TypeByExtension(filePath)
	if mimeType == "" {
		buffer := make([]byte, 512)
		n, err := file.Read(buffer)
		if err != nil && err.Error() != "EOF" {
			return "", err
		}
		mimeType = http.DetectContentType(buffer[:n])
	}

	return mimeType, nil
}
