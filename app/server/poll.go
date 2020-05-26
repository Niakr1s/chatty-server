package server

import (
	"encoding/json"
	"net/http"

	"github.com/niakr1s/chatty-server/app/constants"
)

// Poll ...
func (s *Server) Poll(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(constants.CtxUserNameKey).(string)

	event := <-s.pool.GetUserChan(username)

	json.NewEncoder(w).Encode(event)
}
