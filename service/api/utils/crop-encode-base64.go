package utils

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/jpeg"
	"os"

	"github.com/nfnt/resize"
)

// CropAndEncodeBase64 crops an image to a square, resizes it, and returns it as Base64
func CropAndEncodeBase64(imagePath string, size uint) (string, error) {
	// Open the image file
	file, err := os.Open(imagePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	// Get the dimensions
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Determine the crop area (centered square)
	cropSize := min(width, height)
	x0 := (width - cropSize) / 2
	y0 := (height - cropSize) / 2

	// Crop the image
	croppedImg := img.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(image.Rect(x0, y0, x0+cropSize, y0+cropSize))

	// Resize the cropped image
	resizedImg := resize.Resize(size, size, croppedImg, resize.Lanczos3)

	// Encode the resized image to Base64
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, resizedImg, &jpeg.Options{Quality: 90}); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// Helper function to get the minimum of two numbers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
