package api

import (
	"fmt"
	"image"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"wasaText/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

// rt.router.PUT("/profiles/:userID/groups/:groupID/g_photo", rt.wrap(rt.setGroupPhoto, true))
func (rt *_router) setGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userID, err := strconv.Atoi(ps.ByName("userID"))
	if err != nil {
		BadRequest(w, err, ctx, "Invalid userID")
		return
	}

	// Check if the user is authorized
	if userID != ctx.UserId {
		Forbidden(w, err, ctx, "Unauthorized")
		return
	}

	// Take the group id from the endpoint
	groupID, err := strconv.Atoi(ps.ByName("groupID"))
	if err != nil {
		BadRequest(w, err, ctx, "Invalid groupID")
		return
	}

	// Parse the multipart form
	err = r.ParseMultipartForm(5 << 20) // maxMemory 5MB
	if err != nil {
		http.Error(w, "Bad Request "+err.Error(), http.StatusBadRequest)
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
		log.Fatalf("failed to seek file: %v", err)
	}

	// Define storage path
	uploadPath := fmt.Sprintf("./storage/%d/media/", groupID)
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
	err = rt.db.UpdateGroupPhotoPath(groupID, newFilePath)
	if err != nil {
		InternalServerError(w, err, ctx)
		return
	}

	w.WriteHeader(http.StatusOK)
}
