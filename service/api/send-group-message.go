package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"wasaText/service/api/reqcontext"
	"wasaText/service/api/utils"
	"wasaText/service/structs"

	"github.com/julienschmidt/httprouter"
)

// 	rt.router.POST("/profiles/:userID/conversations/:groupID/messages", rt.wrap(rt.sendGroupMessage, true))

func (rt *_router) sendGroupMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// parse the sender of the message
	sendID, err := strconv.Atoi(ps.ByName("userID"))
	if err != nil {
		BadRequest(w, err, ctx, "Invalid userID")
		return
	}
	// check sender's auth
	if sendID != ctx.UserId {
		Forbidden(w, err, ctx, "Unauthorized")
		return
	}
	// parse the receiver's id
	groupID, err := strconv.Atoi(ps.ByName("groupID"))
	if err != nil {
		BadRequest(w, err, ctx, "Invalid groupID")
		return
	}

	// Ensure it's a POST request
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseMultipartForm(5 << 20); err != nil { // 5MB max
		BadRequest(w, err, ctx, "Failed to parse form")
		return
	}

	content := r.FormValue("content") // Get text message
	// convert directly to int
	msgType, err := strconv.Atoi(r.FormValue("isPhoto"))
	if err != nil {
		BadRequest(w, err, ctx, "Invalid type value")
		return
	}

	var base64Image string

	// If it's an image (msgType == 1)
	if msgType == 1 {
		// Read the uploaded image file
		file, _, err := r.FormFile("image")
		if err != nil {
			BadRequest(w, err, ctx, "Failed to read uploaded image")
			return
		}
		defer file.Close()

		// Convert the image file to a base64 string
		base64Image, err = utils.EncodeImageToBase64(file)
		if err != nil {
			InternalServerError(w, err, ctx)
			return
		}

		// Set content as base64 encoded image
		content = base64Image
	}

	// Save message to database
	newMessage := structs.Message{
		ConvoID:  groupID,
		SenderID: sendID,
		Content:  content,
		IsPhoto:  msgType, // 0 for text, 1 for image
	}
	newMessage.IsForwarded = 0
	// insert message inside db
	msgID, timestamp, err := rt.db.InsertGroupMessage(newMessage)
	if err != nil {
		InternalServerError(w, err, ctx)
		return
	}
	newMessage.MsgID = msgID
	newMessage.Timestamp = timestamp
	newMessage.CheckSent = true

	// response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(newMessage); err != nil {
		ctx.Logger.WithError(err).Error("Error encoding response")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
