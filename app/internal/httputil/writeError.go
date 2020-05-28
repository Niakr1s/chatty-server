package httputil

import (
	"encoding/json"
	"net/http"

	"github.com/niakr1s/chatty-server/app/er"
)

// WriteError ...
func WriteError(w http.ResponseWriter, err error, code int) {
	w.WriteHeader(code)

	jsonErr := struct {
		What string `json:"error"`
	}{err.Error()}

	json.NewEncoder(w).Encode(jsonErr)
}

// WriteSessionError ...
func WriteSessionError(w http.ResponseWriter) {
	WriteError(w, er.ErrSession, http.StatusInternalServerError)
}
