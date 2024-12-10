package api

import (
	"net/http"
	"wasaText/service/api/reqcontext"
)

func InternalServerError(w http.ResponseWriter, err error, ctx reqcontext.RequestContext) {
	ctx.Logger.WithError(err).Error("failed to connect to the database")
	w.WriteHeader(http.StatusInternalServerError)
}

func BadRequest(w http.ResponseWriter, err error, ctx reqcontext.RequestContext, message string) {
	if err != nil {
		http.Error(w, message+": "+err.Error(), http.StatusBadRequest)
	} else {
		http.Error(w, message, http.StatusBadRequest)
	}
}

func Forbidden(w http.ResponseWriter, err error, ctx reqcontext.RequestContext, message string) {
	if err != nil {
		http.Error(w, message+": "+err.Error(), http.StatusForbidden)
	} else {
		http.Error(w, message, http.StatusForbidden)
	}
}
