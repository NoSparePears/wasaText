package api

import (
	"encoding/json"
	"net/http"
	"regexp"

	"wasaText/service/api/reqcontext"
	"wasaText/service/structs"

	"github.com/julienschmidt/httprouter"
)

// rt.router.GET("/profiles", rt.wrap(rt.searchUser, true))
func (rt *_router) searchUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Get the search query from the request
	query_search := r.URL.Query().Get("username")
	validQuerySearch := regexp.MustCompile(`^[a-z0-9]{1,13}$`)
	if !validQuerySearch.MatchString(query_search) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// check username in database
	dbUsers, err := rt.db.SearchUsers(query_search)
	if err != nil {
		ctx.Logger.Error("Error searching users ", err)
		http.Error(w, "Error searching users", http.StatusInternalServerError)
		return
	}
	users := make([]structs.User, len(dbUsers))
	copy(users, dbUsers)

	// set response header for json content
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// encode user in json
	if err = json.NewEncoder(w).Encode(users); err != nil {
		ctx.Logger.WithError(err).Error("Error encoding response")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
