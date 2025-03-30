package api

import (
	"encoding/json"
	"io"
	"log"
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
	// parse the message's id
	msgID, err := strconv.Atoi(ps.ByName("msgID"))
	if err != nil {
		BadRequest(w, err, ctx, "Invalid msgID")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading body:", err)
		http.Error(w, "Failed to read body", http.StatusInternalServerError)
		return
	}
	log.Println("Request body:", string(body))

	// Now decode it into the struct
	var comment structs.Comment
	err = json.Unmarshal(body, &comment)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	comment.MsgID = msgID
	comment.SenderID = userID

	// check if the user already commented the message
	err = rt.db.CheckComment(msgID, userID)
	if err != nil {

		err = rt.db.DeleteComment(msgID, userID)
		if err != nil {
			InternalServerError(w, err, ctx)
			return
		}
	}

	dbComm, err := rt.db.InsertComment(comment)
	if err != nil {
		InternalServerError(w, err, ctx)
		return
	}
	// response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(dbComm); err != nil {
		ctx.Logger.WithError(err).Error("Error encoding response")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

}
