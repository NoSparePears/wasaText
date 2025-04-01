package api

import (
	"encoding/json"
	"errors"
	"io"
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
		BadRequest(w, err, ctx, "Invalid userID format")
		return
	}
	if userID != ctx.UserId {
		Forbidden(w, errors.New("unauthorized"), ctx, "You are not allowed to create a group for this user")
		return
	}

	// Read and parse request body
	defer r.Body.Close() // Ensure body is closed even if parsing fails

	body, err := io.ReadAll(r.Body)
	if err != nil {
		BadRequest(w, err, ctx, "Failed to read request body")
		return
	}

	var input struct {
		Name    string         `json:"name"`
		Members []structs.User `json:"members"`
	}

	if err = json.Unmarshal(body, &input); err != nil {
		BadRequest(w, err, ctx, "Invalid JSON format")
		return
	}

	if input.Name == "" {
		BadRequest(w, errors.New("missing group name"), ctx, "Group name is required")
		return
	}

	// Create group in database
	ctx.Logger.Info("Creating group in database")
	group, err := rt.db.CreateGroup(input.Name, userID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to create group in database")
		InternalServerError(w, err, ctx)
		return
	}
	ctx.Logger.Infof("Group created with ID: %d", group.GlobalConvoID)

	// Retrieve creator details
	ctx.Logger.Info("Fetching creator details")
	creator, err := rt.db.SearchUserID(userID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to fetch creator details")
		InternalServerError(w, err, ctx)
		return
	}
	group.Members = append(group.Members, creator)
	group.GroupName = input.Name

	// Add selected users to group
	for _, member := range input.Members {
		ctx.Logger.Infof("Adding user ID %d to group", member.ID)
		if err := rt.db.AddToGroup(member.ID, group.GlobalConvoID); err != nil {
			ctx.Logger.WithError(err).Errorf("Failed to add user ID %d to group", member.ID)
			InternalServerError(w, err, ctx)
			return
		}
		group.Members = append(group.Members, member)
	}

	// Set response headers and send JSON response
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(group); err != nil {
		ctx.Logger.WithError(err).Error("Error encoding response")
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	ctx.Logger.Info("Group created successfully and response sent")
}
