package api

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"

	"wasaText/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) validGroupName(name string) bool {
	str := name
	validName := regexp.MustCompile(`^[A-Za-zÀ-ÖØ-öø-ÿ0-9 ]{3,100}(/[A-Za-z])?$`)
	return validName.MatchString(str)
}

// rt.router.PUT("/profiles/:userID/groups/:groupID/g_name", rt.wrap(rt.setGroupName, true))
func (rt *_router) setGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userID, err := strconv.Atoi(ps.ByName("userID"))
	if err != nil {
		BadRequest(w, err, ctx, "Invalid userID")
		return
	}
	if userID != ctx.UserId {
		Forbidden(w, err, ctx, "Unauthorized")
		return
	}
	groupID, err := strconv.Atoi(ps.ByName("groupID"))
	if err != nil {
		BadRequest(w, err, ctx, "Invalid groupID")
		return
	}

	var name string

	if err := json.NewDecoder(r.Body).Decode(&name); err != nil {
		http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
		return
	}

	if !rt.validGroupName(name) {
		http.Error(w, "Invalid name for the group", http.StatusBadRequest)
		return
	}

	err = rt.db.SetGroupName(groupID, name)
	if err != nil {
		InternalServerError(w, err, ctx)
		return
	}

	// set response header for json content
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// encode group in json
	if err = json.NewEncoder(w).Encode("Group's name changed correctly."); err != nil {
		ctx.Logger.WithError(err).Error("Error encoding response")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

}
