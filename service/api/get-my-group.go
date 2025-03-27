package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"wasaText/service/api/reqcontext"
	"wasaText/service/structs"

	"github.com/julienschmidt/httprouter"
)

// 	rt.router.GET("/profiles/:userID/groups/:groupID", rt.wrap(rt.getMyGroup, true))

func (rt *_router) getGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	// parse and validate groupID
	groupID, err := strconv.Atoi(ps.ByName("groupID"))
	if err != nil {
		InternalServerError(w, err, ctx)
		return
	}

	// retrieve and validate convo data for the group from db (by groupID)
	dbGroupMembers, err := rt.db.GetGroupMembers(groupID)
	if err != nil {
		InternalServerError(w, err, ctx)
		return
	}

	var members []structs.User

	for _, dbGroupMember := range dbGroupMembers {
		user, err := rt.db.GetUsernameByID(dbGroupMember.ID)
		if err != nil {
			InternalServerError(w, err, ctx)
			return
		}
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
