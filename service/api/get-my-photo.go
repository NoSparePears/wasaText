package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"wasaText/service/api/reqcontext"
	"wasaText/service/api/utils"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getMyPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// parse and validate userID
	userID, err := strconv.Atoi(ps.ByName("userID"))
	if err != nil {
		BadRequest(w, err, ctx, "Invalid userID")
		return
	}

	// check authorization
	if userID != ctx.UserId {
		Forbidden(w, err, ctx, "Unauthorized")
		return
	}

	// Get the photo from the database
	pfpPath, err := rt.db.GetUserPhotoPath(userID)
	if err != nil {
		InternalServerError(w, err, ctx)
		return
	}

	// Crop & encode the profile picture as Base64
	pfpBase64, err := utils.CropAndEncodeBase64(pfpPath, 200) // Crop to 200x200px
	if err != nil {
		InternalServerError(w, err, ctx)
		return
	}

	// Send the cropped image as JSON
	response := map[string]string{
		"profile_picture": pfpBase64,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
