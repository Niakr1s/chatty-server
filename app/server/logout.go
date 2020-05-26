package server

import (
	"net/http"

	"github.com/niakr1s/chatty-server/app/constants"
	"github.com/niakr1s/chatty-server/app/internal/httputil"
)

// Logout - without password!
func (s *Server) Logout(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(constants.CtxUserNameKey).(string)

	s.dbStore.LoggedDB.Lock()
	defer s.dbStore.LoggedDB.Unlock()

	err := s.dbStore.LoggedDB.Logout(username)
	if err != nil {
		httputil.WriteError(w, err, http.StatusConflict)
		return
	}
}
