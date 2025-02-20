package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"wasaText/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

//	rt.router.POST("/profiles/:userID/conversations/:destID/messages/:msgID", rt.wrap(rt.forwardMessage, true))

func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	userID, err := strconv.Atoi(ps.ByName("userID"))
	if err != nil {
		BadRequest(w, err, ctx, "invalid userID")
		return
	}
	if userID != ctx.UserId {
		Forbidden(w, err, ctx, "Unauthorized")
		return
	}

	destID, err := strconv.Atoi(ps.ByName("destID"))
	if err != nil {
		BadRequest(w, err, ctx, "invalid destID")
		return
	}

	msgID, err := strconv.Atoi(ps.ByName("msgID"))
	if err != nil {
		BadRequest(w, err, ctx, "invalid userID")
		return
	}

	// Parsing del corpo della richiesta
	var forwardRequest struct {
		TargetConvoID int    `json:"targetConvoID"` // Dove inoltrare il messaggio
		CustomContent string `json:"customContent,omitempty"`
	}
	err = json.NewDecoder(r.Body).Decode(&forwardRequest)
	if err != nil || forwardRequest.TargetConvoID == 0 {
		BadRequest(w, err, ctx, "Invalid request body")
		return
	}
	//controllo che la conversazione nella quale voglio inoltrare msg esista, sennò la creo
	convo, err := rt.db.GetConversation(userID, forwardRequest.TargetConvoID)
	if err != nil {
		InternalServerError(w, err, ctx)
	}

	msgToForward, err := rt.db.ComposeMsgToForward(msgID, destID, forwardRequest.CustomContent)
	if err != nil {
		InternalServerError(w, err, ctx)
		return
	}
	msgToForward.SenderID = userID

	finalMsg, err := rt.db.InsertMessage(msgToForward, convo.GlobalConvoID)

	// Risposta JSON con il nuovo messaggio inoltrato
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(finalMsg); err != nil {
		ctx.Logger.WithError(err).Error("Error encoding response")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

}
