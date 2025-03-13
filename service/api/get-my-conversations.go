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

// rt.router.GET("/profiles/:userID/conversations", rt.wrap(rt.getMyConversations, true))

func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

	// Retrieve conversations
	dbConvos, err := rt.db.GetConversations(userID)
	if err != nil {
		InternalServerError(w, err, ctx)
		return
	}

	type response struct {
		DbConvo  structs.Conversation `json:"conversation"`
		DestUser structs.User         `json:"destUser"`
		LastMsg  structs.Message      `json:"lastMessage"`
	}

	var convos []response

	for _, dbConvo := range dbConvos {
		destID := dbConvo.DestUserID
		destUser, err := rt.db.GetUsernameByID(destID)
		if err != nil {
			InternalServerError(w, err, ctx)
			return
		}
		destUser.ID = destID

		// Get last message (only if it exists)
		lastMsg := structs.Message{Content: "No messages yet"}
		if dbConvo.LastMsgID != 0 {
			lastMsg, err = rt.db.GetLastMessage(dbConvo.GlobalConvoID, dbConvo.LastMsgID)
			if err != nil {
				InternalServerError(w, err, ctx)
				return
			}
		}

		convos = append(convos, response{
			DbConvo:  dbConvo,
			DestUser: destUser,
			LastMsg:  lastMsg,
		})
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(convos); err != nil {
		InternalServerError(w, err, ctx)
	}
}
