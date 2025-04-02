package api

import (
	"fmt"
	"image"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	_ "image/jpeg" // Required for decoding JPEG files
	_ "image/png"  // Required for decoding PNG files

	"wasaText/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

// Route: PUT /profiles/:userID/photo
func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Parse userID
	userID, err := strconv.Atoi(ps.ByName("userID"))
	if err != nil {
		http.Error(w, "Invalid user ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Check if user is authorized
	if userID != ctx.UserId {
		Forbidden(w, nil, ctx, "Unauthorized")
		return
	}

	// Parse multipart form (limit: 5MB)
	err = r.ParseMultipartForm(5 << 20)
	if err != nil {
		http.Error(w, "Error parsing form: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve the uploaded file
	file, _, err := r.FormFile("profile_picture") // No need for `header`
	if err != nil {
		BadRequest(w, err, ctx, "Invalid file upload")
		return
	}
	defer file.Close()

	// Validate image format
	_, format, err := image.Decode(file) // Only checking format
	if err != nil || (format != "jpeg" && format != "png") {
		BadRequest(w, err, ctx, "Invalid image file (must be JPG or PNG)")
		return
	}

	// Reset file reader (needed after decoding)
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		BadRequest(w, err, ctx, fmt.Sprintf("failed to seek file: %v", err))
		return
	}

	// Define storage path
	uploadPath := fmt.Sprintf("./storage/%d/media/", userID)
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		InternalServerError(w, err, ctx)
		return
	}

	// Always save as JPG (convert PNG if necessary)
	newFilePath := filepath.Join(uploadPath, "profile.jpg")

	// Create output file
	outFile, err := os.Create(newFilePath)
	if err != nil {
		InternalServerError(w, err, ctx)
		return
	}
	defer outFile.Close()

	// Copy the file to storage
	_, err = io.Copy(outFile, file)
	if err != nil {
		InternalServerError(w, err, ctx)
		return
	}

	// Update database with new profile picture path
	err = rt.db.UpdateUserPhotoPath(userID, newFilePath)
	if err != nil {
		InternalServerError(w, err, ctx)
		return
	}

	w.WriteHeader(http.StatusOK)

}
