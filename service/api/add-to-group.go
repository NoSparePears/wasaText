package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"wasaText/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

// rt.router.PUT("/profiles/:userID/groups/:groupID/members", rt.wrap(rt.addToGroup, true))
func (rt *_router) addToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// parse user id
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

	// Define a struct to match JSON body
	var requestBody struct {
		MemberID int `json:"memberID"`
	}

	// Decode JSON request body into the struct
	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		BadRequest(w, err, ctx, "Error decoding JSON")
		return
	}

	// Get the parsed member ID
	memberID := requestBody.MemberID

	// Ensure memberID is valid
	if memberID <= 0 {
		BadRequest(w, err, ctx, "Invalid memberID")
		return
	}

	err = rt.db.AddToGroup(memberID, groupID)
	if err != nil {
		InternalServerError(w, err, ctx)
		return
	}

	// set response header for json content

	w.WriteHeader(http.StatusOK)

}
