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

// rt.router.PUT("/profiles/:userID/conversations/:destID/messages/:msgID/comments", rt.wrap(rt.commentMessage, true))
func (rt *_router) commentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userID, err := strconv.Atoi(ps.ByName("userID"))
	if err != nil {
		BadRequest(w, err, ctx, "invalid userID")
		return
	}
	if userID != ctx.UserId {
		Forbidden(w, err, ctx, "Unauthorized")
		return
	}
	//parse the message's id
	msgID, err := strconv.Atoi(ps.ByName("msgID"))
	if err != nil {
		BadRequest(w, err, ctx, "Invalid msgID")
		return
	}

	var comment structs.Comment
	//read request body
	err = json.NewDecoder(r.Body).Decode(&comment.Emoji)
	if err != nil {
		BadRequest(w, err, ctx, "Error decoding JSON")
		return
	}
	//check provided body for the comment
	if comment.Emoji == "" {
		BadRequest(w, errors.New("string is required"), ctx, "Missing message body")
		return
	}

	comment.MsgID = msgID
	comment.SenderID = userID

	dbComm, err := rt.db.InsertComment(comment)
	if err != nil {
		InternalServerError(w, err, ctx)
		return
	}
	//response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(dbComm); err != nil {
		ctx.Logger.WithError(err).Error("Error encoding response")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

}
