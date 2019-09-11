package io

import (
	"net/http"
	"os"
)

func GetFileContentType(file *os.File) (string, error) {
	buffer := make([]byte, 512) // we only need the first 512 bytes of the file to be able to determine its file type
	_, err := file.Read(buffer)

	if err != nil {
		return "", err
	}

	contentType := http.DetectContentType(buffer)
	return contentType, nil

}
