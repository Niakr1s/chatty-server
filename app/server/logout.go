package server

import (
	"net/http"

	"github.com/niakr1s/chatty-server/app/config"
	"github.com/niakr1s/chatty-server/app/server/httputil"
)

// Logout - without password!
func (s *Server) Logout(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(config.CtxUserNameKey).(string)

	s.Store.LoggedDB.Lock()
	defer s.Store.LoggedDB.Unlock()

	err := s.Store.LoggedDB.Logout(username)
	if err != nil {
		httputil.WriteError(w, err, http.StatusConflict)
		return
	}
}
