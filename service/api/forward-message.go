package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"wasaText/service/api/reqcontext"
	"wasaText/service/structs"

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

	msgID, err := strconv.Atoi(ps.ByName("msgID"))
	if err != nil {
		BadRequest(w, err, ctx, "invalid userID")
		return
	}

	// Ensure it's a POST request
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// parse  {destID: destUser.id} from the body
	var requestBody struct {
		DestID int `json:"destID"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		BadRequest(w, err, ctx, "Invalid request body")
		return
	}
	destID := requestBody.DestID
	if destID == 0 {
		BadRequest(w, err, ctx, "Invalid destID")
		return
	}

	content, isPhoto, err := rt.db.GetMessageContent(msgID)
	if err != nil {
		InternalServerError(w, err, ctx)
		return
	}

	var msgToForward structs.Message

	msgToForward.SenderID = userID
	msgToForward.Content = content
	msgToForward.IsPhoto = isPhoto
	msgToForward.IsForwarded = 1 // Set IsForwarded to 1

	convoID, err := rt.db.GetConvoID(userID, destID)
	if err != nil {
		InternalServerError(w, err, ctx)
		return
	}

	msgToForward.ConvoID = convoID

	id, timestamp, err := rt.db.InsertMessage(msgToForward)
	if err != nil {
		InternalServerError(w, err, ctx)
		return
	}

	msgToForward.MsgID = id
	msgToForward.Timestamp = timestamp
	msgToForward.CheckSent = true

	// Risposta JSON con il nuovo messaggio inoltrato
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(msgToForward); err != nil {
		ctx.Logger.WithError(err).Error("Error encoding response")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

}
