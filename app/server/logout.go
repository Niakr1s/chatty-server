package server

import (
	"net/http"

	"github.com/niakr1s/chatty-server/app/internal/httputil"
	"github.com/niakr1s/chatty-server/app/internal/sess"
)

// Logout - without password!
func (s *Server) Logout(w http.ResponseWriter, r *http.Request) {
	username := sess.GetUserNameFromCtx(r.Context())

	s.dbStore.LoggedDB.Lock()
	defer s.dbStore.LoggedDB.Unlock()

	err := s.dbStore.LoggedDB.Logout(username)
	if err != nil {
		httputil.WriteError(w, err, http.StatusConflict)
		return
	}
}
