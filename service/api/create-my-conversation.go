package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"wasaText/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)
// rt.router.PUT("/profiles/:userID/conversations/:destID", rt.wrap(rt.createConversation, true))
func (rt *_router) createConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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
	
	var convo Conversation 
	convo.SendID = userID
	convo.RecID = destID 

	//return convo object with his identificator
	dbConvo, err := rt.db.CreateConversation(convo.ToDatabase())
	if err != nil {
		return convo, err 
	}

	//checks if the convo has been succesfully added to the db
	err = convo.FromDatabase(dbConvo)
	if err != nil {
		return convo, err 
	}

	return convo, nil
}
