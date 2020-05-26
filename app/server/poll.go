package server

import (
	"encoding/json"
	"net/http"

	"github.com/niakr1s/chatty-server/app/constants"
	"github.com/niakr1s/chatty-server/app/internal/httputil"
)

// Poll ...
func (s *Server) Poll(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(constants.CtxUserNameKey).(string)

	event := <-s.pool.GetUserChan(username)

	if err := json.NewEncoder(w).Encode(event); err != nil {
		httputil.WriteError(w, err, http.StatusInternalServerError)
		return
	}
}
