package util

import (
	"encoding/json"
	"net/http"
)

// WriteError ...
func WriteError(w http.ResponseWriter, err error, code int) {
	w.WriteHeader(code)

	jsonErr := struct {
		What string `json:"error"`
	}{err.Error()}

	json.NewEncoder(w).Encode(jsonErr)
}
