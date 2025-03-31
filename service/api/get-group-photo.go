package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"wasaText/service/api/reqcontext"
	"wasaText/service/api/utils"

	"github.com/julienschmidt/httprouter"
)

// 	rt.router.GET("/profiles/:userID/groups/:groupID/g_photo", rt.wrap(rt.getGroupPhoto, true))

func (rt *_router) getGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// parse and validate userID
	userID, err := strconv.Atoi(ps.ByName("userID"))
	if err != nil {
		BadRequest(w, err, ctx, "Invalid userID")
		return
	}

	if userID != ctx.UserId {
		Forbidden(w, err, ctx, "Unauthorized")
		return
	}

	groupID, err := strconv.Atoi(ps.ByName("groupID"))
	if err != nil {
		BadRequest(w, err, ctx, "Invalid groupID")
		return
	}

	// Get the photo from the database
	pfpPath, err := rt.db.GetGroupPhotoPath(groupID)
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

	// encode user in json
	if err = json.NewEncoder(w).Encode(response); err != nil {
		ctx.Logger.WithError(err).Error("Error encoding response")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
