package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"wasaText/service/api/reqcontext"
	"wasaText/service/structs"

	"github.com/julienschmidt/httprouter"
)

// rt.router.POST("/profiles/:userID/groups", rt.wrap(rt.createGroup, true))
func (rt *_router) createGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userID, err := strconv.Atoi(ps.ByName("userID"))
	if err != nil {
		BadRequest(w, err, ctx, "Invalid userID")
		return
	}
	if userID != ctx.UserId {
		Forbidden(w, err, ctx, "Unauthorized")
		return
	}

	// Ensure it's a POST request
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	type groupRequest struct {
		Name    string         `json:"name"`
		Photo   string         `json:"photo"`
		Members []structs.User `json:"members"`
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Parse JSON
	var input groupRequest
	err = json.Unmarshal(body, &input)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	// check provided string
	if input.Name == "" {
		BadRequest(w, errors.New("string is required"), ctx, "Missing group's name")
		return
	}

	// create globalConvoTable con isGroup = 1 e groupname = input da frontend, inoltre crea la groupMemberTable con userID = userID
	group, err := rt.db.CreateGroup(input.Name, userID)
	if err != nil {
		InternalServerError(w, err, ctx)
		return
	}
	creator, err := rt.db.SearchUserID(userID)
	if err != nil {
		InternalServerError(w, err, ctx)
		return
	}
	group.Members = append(group.Members, creator)
	group.GroupPropic64 = input.Photo

	// create membertable per ogni selectedUser da frontend
	for _, member := range input.Members {
		err = rt.db.AddToGroup(member.ID, group.GlobalConvoID)
		if err != nil {
			InternalServerError(w, err, ctx)
			return
		}
		group.Members = append(group.Members, member)
	}

	type response struct {
		Group structs.Group `json:"group"`
	}
	resp := response{
		Group: group,
	}
	// SCEGLI SE CREARE LA CONVO TABLE PER OGNI UTENTE ALL INTERNO DI ADDTOGROUP (LATO DB) OPPURE FAI LE CHIAMATE QUI
	// set response header for json content
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// encode user in json
	if err = json.NewEncoder(w).Encode(resp); err != nil {
		ctx.Logger.WithError(err).Error("Error encoding response")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

}
