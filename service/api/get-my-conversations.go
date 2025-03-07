package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"wasaText/service/api/reqcontext"
	"wasaText/service/structs"

	"github.com/julienschmidt/httprouter"
)

// rt.router.GET("/profiles/:userID/conversations", rt.wrap(rt.getMyConversations, true))

func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	//parse and validate userID
	userID, err := strconv.Atoi(ps.ByName("userID"))
	if err != nil {
		BadRequest(w, err, ctx, "Invalid userID")
		return
	}

	//check user authorization
	if userID != ctx.UserId {
		Forbidden(w, err, ctx, "Unauthorized")
		return
	}

	var dbConvos []structs.Conversation
	//retrieve and validate convos data from db
	dbConvos, err = rt.db.GetConversations(userID)
	if err != nil {
		InternalServerError(w, err, ctx)
		return
	}
	type response struct {
		DbConvo  structs.Conversation `json:"conversation"`
		DestUser structs.User         `json:"destUser"`
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

		convo := response{
			DbConvo:  dbConvo,
			DestUser: destUser,
		}
		convos = append(convos, convo)
	}

	//response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(convos); err != nil {
		ctx.Logger.WithError(err).Error("Error encoding response")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
