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

func (rt *_router) getMyGroups(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Parse and validate userID
	userID, err := strconv.Atoi(ps.ByName("userID"))
	if err != nil {
		BadRequest(w, err, ctx, "Invalid userID")
		return
	}

	// Check user authorization
	if userID != ctx.UserId {
		Forbidden(w, errors.New("user unauthorized"), ctx, "Unauthorized")
		return
	}

	// Retrieve groups
	dbGroups, err := rt.db.GetGroups(userID)
	if err != nil {
		InternalServerError(w, err, ctx)
		return
	}

	type response struct {
		Group   structs.Group   `json:"group"`
		LastMsg structs.Message `json:"lastMessage"`
	}

	var groups []response

	for _, dbGroup := range dbGroups {
		// Get last message (only if it exists)
		lastMsg := structs.Message{Content: "No messages yet"}
		if dbGroup.LastMsgID != 0 {
			lastMsg, err = rt.db.GetLastMessage(dbGroup.GlobalConvoID, dbGroup.LastMsgID)
			if err != nil {
				InternalServerError(w, err, ctx)
				return
			}
		}

		groups = append(groups, response{
			Group:   dbGroup,
			LastMsg: lastMsg,
		})
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(groups); err != nil {
		InternalServerError(w, err, ctx)
	}
}
