package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"wasaText/service/api/reqcontext"
	"wasaText/service/structs"

	"github.com/julienschmidt/httprouter"
)

// rt.router.DELETE("/profiles/:userID/conversations/:destID", rt.wrap(rt.deleteConversation, true))
func (rt *_router) deleteConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	//parse requesting user ID
	reqUserID, err := strconv.Atoi(ps.ByName("userID"))
	if err != nil {
		BadRequest(w, err, ctx, "Invalid userID")
		return
	}
	//check auth
	if reqUserID != ctx.UserId {
		Forbidden(w, err, ctx, "Unauthorized")
		return
	}
	//parse destID
	destID, err := strconv.Atoi((ps.ByName("destID")))
	if err != nil {
		BadRequest(w, err, ctx, "Invalid destID")
		return
	}
	var convos []structs.Conversation

	convos, err = rt.db.DeleteConversation(reqUserID, destID)
	if err != nil {
		InternalServerError(w, err, ctx)
		return
	}

	//set responde header for json content
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	//encode convos in json
	if err = json.NewEncoder(w).Encode(convos); err != nil {
		ctx.Logger.WithError(err).Error("Error encoding response")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

}
