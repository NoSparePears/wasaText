package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"wasaText/service/api/reqcontext"
	"wasaText/service/structs"

	"github.com/julienschmidt/httprouter"
)

// rt.router.POST("/profiles/:userID/groups", rt.wrap(rt.createGroup, true))
func (rt *_router) createGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userID, err := strconv.Atoi(ps.ByName("userID"))
	if err != nil {
		BadRequest(w, err, ctx, "Invalid userID")
		return
	}
	if userID != ctx.UserId {
		Forbidden(w, err, ctx, "Unauthorized")
		return
	}

	var groupRequest struct {
		Name    string         `json:"groupName"`
		Members []structs.User `json:"members"`
	}
	err = json.NewDecoder(r.Body).Decode(&groupRequest)
	if err != nil {
		BadRequest(w, err, ctx, "Error decoding JSON")
		return
	}
	// check provided string
	if groupRequest.Name == "" {
		BadRequest(w, errors.New("string is required"), ctx, "Missing group's name")
		return
	}

	group, err := rt.db.CreateGroup(groupRequest.Name, userID)
	if err != nil {
		InternalServerError(w, err, ctx)
		return
	}
	for _, member := range groupRequest.Members {
		err = rt.db.AddToGroup(member.ID, group.GlobalConvoID)
		if err != nil {
			InternalServerError(w, err, ctx)
			return
		}
		group.Members = append(group.Members, member)
	}

	// set response header for json content
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// encode user in json
	if err = json.NewEncoder(w).Encode(group); err != nil {
		ctx.Logger.WithError(err).Error("Error encoding response")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

}
