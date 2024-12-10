package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"wasaText/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// rt.router.GET("/profiles/:userID/conversations/:destID", rt.wrap(rt.getConversation, true))

func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	
	//parse and validate userID
	userID, err := strconv.Atoi(ps.ByName("userID"))
	if err != nil {
		BadRequest(w, err, ctx, "Invalid userID")
		return
	}

	//check authorization
	if userID != ctx.ID {
		Forbidden(w, err, ctx, "Unauthorized")
		return 
	}

	//parse and validate destID
	destID, err := strconv.Atoi (ps.ByName("destID"))
	if err != nil {
		InternalServerError(w, err, ctx)
		return 
	}

	//retrieve and validate convo data from db
	dbConvo, err := rt.db.GetConversation(userID, destID)
	if err != nil {
		InternalServerError(w, err, ctx)
		return 
	}

	//map convo model from db to api model
	var convo Conversation
	convo, err := rt.FromDatabase(dbconvo)
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
