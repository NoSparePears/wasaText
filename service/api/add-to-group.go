package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"wasaText/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

// rt.router.PUT("/profiles/:userID/groups/:groupID/members", rt.wrap(rt.addToGroup, true))
func (rt *_router) addToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// parse user id
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

	var memberID int

	err = json.NewDecoder(r.Body).Decode(&memberID)
	if err != nil {
		BadRequest(w, err, ctx, "Error decoding JSON")
		return
	}

	err = rt.db.AddToGroup(memberID, groupID)
	if err != nil {
		InternalServerError(w, err, ctx)
		return
	}

	members, err := rt.db.GetGroupMembers(groupID)
	if err != nil {
		InternalServerError(w, err, ctx)
		return
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
