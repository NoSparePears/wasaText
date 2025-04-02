package api

import (
	"encoding/json"
	"net/http"
	"wasaText/service/api/reqcontext"
	"wasaText/service/structs"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var user structs.User

	// read request body
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		BadRequest(w, err, ctx, "Error decoding JSON")
		return
	}
	// Check if the patter of the username respect the regex
	if !user.IsValid() {
		BadRequest(w, err, ctx, "Invalid username")
		return
	}
	exists, err := rt.db.CheckUsername(user.Username)
	if err != nil {
		InternalServerError(w, err, ctx)
		return
	}
	if !exists {
		dbUser, err := rt.db.CreateUser(user.Username)
		if err != nil {
			InternalServerError(w, err, ctx)
			return
		}
		user = dbUser
	} else {
		dbUser, err := rt.db.SearchUser(user.Username)
		if err != nil {
			InternalServerError(w, err, ctx)
			return
		}
		user = dbUser
	}

	// Response
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		InternalServerError(w, err, ctx)
		return
	}

}
