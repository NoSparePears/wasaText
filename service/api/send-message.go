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

//rt.router.POST("/profiles/:userID/conversations/:destID/messages", rt.wrap(rt.sendMessage, true))

func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	//parse the sender of the message
	sendID, err := strconv.Atoi(ps.ByName("userID"))
	if err != nil {
		BadRequest(w, err, ctx, "Invalid userID")
		return
	}
	//check sender's auth
	if sendID != ctx.UserId {
		Forbidden(w, err, ctx, "Unauthorized")
		return
	}
	//parse the receiver's id
	recID, err := strconv.Atoi(ps.ByName("destID"))
	if err != nil {
		BadRequest(w, err, ctx, "Invalid destID")
		return
	}
	//check convo in db
	convo, err := rt.db.GetConversation(sendID, recID)

	var msg structs.Message
	//read request body
	err = json.NewDecoder(r.Body).Decode(&msg.Content)
	if err != nil {
		BadRequest(w, err, ctx, "Error decoding JSON")
		return
	}
	//check provided body for the message
	if msg.Content == "" {
		BadRequest(w, errors.New("string is required"), ctx, "Missing message body")
		return
	}

	msg.ConvoID = convo.GlobalConvoID
	msg.SenderID = sendID
	//insert message inside db
	dbMsg, err := rt.db.InsertMessage(msg, recID)
	if err != nil {
		InternalServerError(w, err, ctx)
		return
	}
	//response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(dbMsg); err != nil {
		ctx.Logger.WithError(err).Error("Error encoding response")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
