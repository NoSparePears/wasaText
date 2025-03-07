package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"wasaText/service/api/reqcontext"
	"wasaText/service/structs"

	"github.com/julienschmidt/httprouter"
)

// rt.router.GET("/profiles/:userID/conversations/:destID", rt.wrap(rt.getConversation, true))

func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// parse and validate userID
	userID, err := strconv.Atoi(ps.ByName("userID"))
	if err != nil {
		BadRequest(w, err, ctx, "Invalid userID")
		return
	}

	// check authorization
	if userID != ctx.UserId {
		Forbidden(w, err, ctx, "Unauthorized")
		return
	}

	// parse and validate destID
	destID, err := strconv.Atoi(ps.ByName("destID"))
	if err != nil {
		InternalServerError(w, err, ctx)
		return
	}

	// retrieve and validate convo data from db
	dbConvo, err := rt.db.GetConversation(userID, destID)
	if err != nil {
		InternalServerError(w, err, ctx)
		return
	}

	destUser, err := rt.db.GetUsernameByID(destID)
	if err != nil {
		InternalServerError(w, err, ctx)
		return
	}
	destUser.ID = destID
	/* IMPLEMENTA PER RICAVARE FOTO PROFILP
	destUser.UserPropic64 = */

	type response struct {
		Convo    structs.Conversation `json:"conversation"`
		DestUser structs.User         `json:"destUser"`
	}
	resp := response{
		Convo:    dbConvo,
		DestUser: destUser,
	}
	// response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(resp); err != nil {
		ctx.Logger.WithError(err).Error("Error encoding response")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
