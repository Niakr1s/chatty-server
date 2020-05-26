package server

import (
	"net/http"

	"github.com/niakr1s/chatty-server/app/config"
	"github.com/niakr1s/chatty-server/app/server/httputil"
)

// Logout - without password!
func (s *Server) Logout(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(config.CtxUserNameKey).(string)

	s.store.LoggedDB.Lock()
	defer s.store.LoggedDB.Unlock()

	err := s.store.LoggedDB.Logout(username)
	if err != nil {
		httputil.WriteError(w, err, http.StatusConflict)
		return
	}
}
