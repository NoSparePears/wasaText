package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"wasaText/service/api/reqcontext"
	"wasaText/service/structs"

	"github.com/julienschmidt/httprouter"
)

// rt.router.GET("/profiles/:userID/groups/:groupID/messages", rt.wrap(rt.getGroupMessages, true))
func (rt *_router) getGroupMessages(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// parse and validate userID
	userID, err := strconv.Atoi(ps.ByName("userID"))
	if err != nil {
		BadRequest(w, err, ctx, "Invalid userID")
		return
	}

	// check user authorization
	if userID != ctx.UserId {
		Forbidden(w, err, ctx, "Unauthorized")
		return
	}

	groupID, err := strconv.Atoi(ps.ByName("groupID"))
	if err != nil {
		BadRequest(w, err, ctx, "Invalid groupID")
		return
	}

	var dbMessages []structs.Message
	// retrieve and validate messages data from db
	dbMessages, err = rt.db.GetGroupMessages(userID, groupID)
	if err != nil {
		InternalServerError(w, err, ctx)
		return
	}

	type response struct {
		Message structs.Message `json:"message"`
		User    structs.User    `json:"user"`
	}

	var messages []response

	for _, dbMessage := range dbMessages {
		sendID := dbMessage.SenderID
		user, err := rt.db.GetUsernameByID(sendID)
		if err != nil {
			InternalServerError(w, err, ctx)
			return
		}
		// Simply append each message to the response slice
		messages = append(messages, response{
			Message: dbMessage,
			User:    user,
		})
	}
	// response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(messages); err != nil {
		ctx.Logger.WithError(err).Error("Error encoding response")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
