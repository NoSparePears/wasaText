package utils

import (
	"encoding/base64"
	"io"
	"io/ioutil"
)

// ImageToBase64 converts an image file to a base64 string
func EncodeImageToBase64(file io.Reader) (string, error) {
	// Read the image file into a byte slice
	imageBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	// Convert the byte slice to a base64 string
	base64Str := base64.StdEncoding.EncodeToString(imageBytes)
	return base64Str, nil
}
