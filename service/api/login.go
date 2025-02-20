package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"wasaText/service/api/reqcontext"
	"wasaText/service/structs"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var user structs.User

	//read request body
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		BadRequest(w, err, ctx, "Error decoding JSON")
		return
	}
	//check username
	dbUser, err := rt.db.SearchUser(user.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) && user.IsValid() {
			//create user if it doesn't exists
			dbUser, err = rt.db.CreateUser(user.Username)
			if err != nil {
				InternalServerError(w, err, ctx)
				return
			}
			w.WriteHeader(http.StatusCreated)
		} else {
			InternalServerError(w, err, ctx)
			return
		}

	}
	//Response
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(dbUser.ID); err != nil {
		InternalServerError(w, err, ctx)
		return
	}

}
