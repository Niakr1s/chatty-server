package server

import (
	"encoding/json"
	"net/http"

	"github.com/niakr1s/chatty-server/app/config"
)

// Poll ...
func (s *Server) Poll(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(config.CtxUserNameKey).(string)

	event := <-s.pool.GetUserChan(username)

	json.NewEncoder(w).Encode(event)
}
