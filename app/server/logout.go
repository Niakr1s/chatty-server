package server

import (
	"net/http"

	"github.com/niakr1s/chatty-server/app/internal/httputil"
	"github.com/niakr1s/chatty-server/app/internal/sess"
)

// Logout - without password!
func (s *Server) Logout(w http.ResponseWriter, r *http.Request) {
	username := sess.GetUserNameFromCtx(r.Context())

	err := s.dbStore.LoggedDB.Logout(username)
	if err != nil {
		httputil.WriteError(w, err, http.StatusConflict)
		return
	}
}
