package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"wasaText/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

//rt.router.DELETE("/profiles/:userID/conversations/:destID/messages/:msgID", rt.wrap(rt.deleteMessage, true))

func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	senderID, err := strconv.Atoi(ps.ByName("userID"))
	if err != nil {
		BadRequest(w, err, ctx, "Invalid userID")
		return
	}

	if senderID != ctx.UserId {
		Forbidden(w, err, ctx, "Unauthorized")
		return
	}

	destID, err := strconv.Atoi(ps.ByName("destID"))
	if err != nil {
		BadRequest(w, err, ctx, "Invalid destID")
		return
	}

	msgID, err := strconv.Atoi(ps.ByName("msgID"))
	if err != nil {
		BadRequest(w, err, ctx, "Invalid msgID")
		return
	}

	convo, err := rt.db.GetConversation(senderID, destID)
	if err != nil {
		InternalServerError(w, err, ctx)
		return
	}

	convo, err = rt.db.DeleteMessage(msgID, convo.GlobalConvoID, senderID)
	if err != nil {
		InternalServerError(w, err, ctx)
		return
	}

	//response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(convo); err != nil {
		ctx.Logger.WithError(err).Error("Error encoding response")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

}
