package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"wasaText/service/api/reqcontext"
	"wasaText/service/structs"

	"github.com/julienschmidt/httprouter"
)

// rt.router.GET("/profiles/:userID/conversations/:destID/messages", rt.wrap(rt.getMessages, true))
func (rt *_router) getMyMessages(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	destID, err := strconv.Atoi(ps.ByName("destID"))
	if err != nil {
		BadRequest(w, err, ctx, "Invalid destID")
		return
	}

	var dbMessages []structs.Message
	// retrieve and validate messages data from db
	dbMessages, err = rt.db.GetMessages(userID, destID)
	if err != nil {
		InternalServerError(w, err, ctx)
		return
	}

	type response struct {
		Message  structs.Message   `json:"message"`
		Comments []structs.Comment `json:"comments"`
	}

	var messages []response

	for _, dbMessage := range dbMessages {

		// Get comments for each message
		comments, err := rt.db.GetMessageComments(dbMessage.MsgID, dbMessage.ConvoID)
		if err != nil {
			ctx.Logger.WithError(err).Error("Error getting comments")
			InternalServerError(w, err, ctx)
			return
		}
		if comments == nil {
			comments = []structs.Comment{}
		}

		messages = append(messages, response{
			Message:  dbMessage,
			Comments: comments,
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
