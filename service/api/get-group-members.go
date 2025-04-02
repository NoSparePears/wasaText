package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"wasaText/service/api/reqcontext"
	"wasaText/service/api/utils"
	"wasaText/service/structs"

	"github.com/julienschmidt/httprouter"
)

// rt.router.GET("/profiles/:userID/groups/:groupID/members", rt.wrap(rt.GetGroupMembers, true))
func (rt *_router) GetGroupMembers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	dbMembers, err := rt.db.GetGroupMembers(groupID)
	if err != nil {
		InternalServerError(w, err, ctx)
		return
	}

	var members []structs.User
	for _, dbMember := range dbMembers {
		user, err := rt.db.GetUsernameByID(dbMember.ID)
		if err != nil {
			InternalServerError(w, err, ctx)
			return
		}
		// Get the photo from the database
		pfpPath, err := rt.db.GetUserPhotoPath(dbMember.ID)
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
		user.UserPropic64 = pfpBase64
		members = append(members, user)
	}

	// set response header for json content
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// encode user in json
	if err = json.NewEncoder(w).Encode(members); err != nil {
		ctx.Logger.WithError(err).Error("Error encoding response")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

}
