package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"wasaText/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

// rt.router.DELETE("/profiles/:userID/conversations/:groupID/messages/:msgID", rt.wrap(rt.deleteGroupMessage, true))

func (rt *_router) deleteGroupMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	senderID, err := strconv.Atoi(ps.ByName("userID"))
	if err != nil {
		BadRequest(w, err, ctx, "Invalid userID")
		return
	}

	if senderID != ctx.UserId {
		Forbidden(w, err, ctx, "Unauthorized")
		return
	}

	groupID, err := strconv.Atoi(ps.ByName("groupID"))
	if err != nil {
		BadRequest(w, err, ctx, "Invalid groupID")
		return
	}

	msgID, err := strconv.Atoi(ps.ByName("msgID"))
	if err != nil {
		BadRequest(w, err, ctx, "Invalid msgID")
		return
	}

	ownerID, err := rt.db.GetMsgOwnerID(msgID)
	if err != nil {
		InternalServerError(w, err, ctx)
		return
	}

	if ownerID != senderID {
		Forbidden(w, err, ctx, "Unauthorized")
		return
	}

	err = rt.db.DeleteMessage(msgID, groupID)
	if err != nil {
		InternalServerError(w, err, ctx)
		return
	}

	messages, err := rt.db.GetGroupMessages(groupID, senderID)
	if err != nil {
		InternalServerError(w, err, ctx)
		return
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
