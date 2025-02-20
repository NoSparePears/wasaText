package utils

import (
	"image/jpeg"
	"os"

	"github.com/nfnt/resize"
)

// Function used to save and crop an image
func SaveAndCrop(filename string, w uint, h uint) error {
	file, err := os.Open(filename) // Opena the image file
	if err != nil {
		return err
	}
	defer func() { err = file.Close() }()

	// Decode the image
	img, err := jpeg.Decode(file)
	if err != nil {
		return err
	}

	resizedImg := resize.Resize(w, h, img, resize.NearestNeighbor)
	// Save cropped image
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() { err = out.Close() }()
	if err := jpeg.Encode(out, resizedImg, nil); err != nil {
		return err
	}

	return err
}
