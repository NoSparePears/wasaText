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

// rt.router.GET("/profiles", rt.wrap(rt.searchUser, true))
func (rt *_router) searchUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var reqUser structs.User

	//read request body
	err := json.NewDecoder(r.Body).Decode(&reqUser)
	if err != nil {
		BadRequest(w, err, ctx, "Error decoding JSON")
		return
	}
	//check provided username
	if reqUser.Username == "" {
		BadRequest(w, errors.New("username is required"), ctx, "Missing username")
		return
	}
	//check username in database
	dbUser, err := rt.db.SearchUser(reqUser.Username)
	if err != nil {
		//user not found
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "User does not exist", http.StatusNotFound)
			return
		}
		//any other error
		InternalServerError(w, err, ctx)
		return
	}

	//set response header for json content
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	//encode user in json
	if err = json.NewEncoder(w).Encode(dbUser); err != nil {
		ctx.Logger.WithError(err).Error("Error encoding response")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
