package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"wasaText/service/api/reqcontext"
	"wasaText/service/structs"

	"github.com/julienschmidt/httprouter"
)

// rt.router.POST("/profiles/:userID/conversations/:destID/messages", rt.wrap(rt.sendMessage, true))

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

	// Ensure it's a POST request
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Parse JSON
	var msg structs.Message
	err = json.Unmarshal(body, &msg)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	msg.SenderID = sendID
	id, err := rt.db.GetConvoID(sendID, recID)
	if err != nil {
		InternalServerError(w, err, ctx)
		return
	}
	msg.ConvoID = id

	//insert message inside db
	dbMsg, err := rt.db.InsertMessage(msg, recID)
	if err != nil {
		InternalServerError(w, err, ctx)
		return
	}
	/*
		type response struct {
			Msg structs.Message `json:"message"`
		}

		resp := response{
			Msg: dbMsg,
		}
	*/
	//response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(dbMsg); err != nil {
		ctx.Logger.WithError(err).Error("Error encoding response")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
