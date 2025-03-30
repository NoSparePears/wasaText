package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"wasaText/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

// rt.router.DELETE("/profiles/:userID/conversations/:destID/messages/:msgID/comments/:commID", rt.wrap(rt.uncommmentMessage, true))

func (rt *_router) uncommmentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	sendID, err := strconv.Atoi(ps.ByName("userID"))
	if err != nil {
		BadRequest(w, err, ctx, "Invalid userID")
		return
	}
	// check auth
	if sendID != ctx.UserId {
		Forbidden(w, err, ctx, "Unauthorized")
		return
	}
	// parse message ID
	msgID, err := strconv.Atoi(ps.ByName("msgID"))
	if err != nil {
		BadRequest(w, err, ctx, "Invalid msgID")
		return
	}

	if err := rt.db.DeleteComment(sendID, msgID); err != nil {
		InternalServerError(w, err, ctx)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "plain/text")
	if err := json.NewEncoder(w).Encode("Comment deleted"); err != nil {
		ctx.Logger.WithError(err).Error("can't encode the reponse")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
