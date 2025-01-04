package utils

import (
	"encoding/hex"
	"fmt"
	"math/rand/v2"
	"mime/multipart"
)

func IsImage(file *multipart.FileHeader) (string, bool) {
	allowedTypes := []string{"image/jpg", "image/png", "image/gif"}
	fileType := file.Header.Get("Content-Type")

	for _, allowedType := range allowedTypes {
		if fileType == allowedType {
			return fileType, true
		}
	}

	return "", false
}

func GenerateImageName(lenght int) (string, error) {
	if lenght%2 != 0 {
		return "", fmt.Errorf("length must be an even number to represent hex bytes")
	}

	bytes := make([]byte, lenght/2)
	rand := rand.ChaCha8{}
	_, err := rand.Read(bytes)

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
