package server

import (
	"net/http"
	"time"

	"github.com/niakr1s/chatty-server/app/er"
	"github.com/niakr1s/chatty-server/app/internal/httputil"
	"github.com/niakr1s/chatty-server/app/internal/sess"
)

// KeepAlive ...
func (s *Server) KeepAlive(w http.ResponseWriter, r *http.Request) {
	username := sess.GetUserNameFromCtx(r.Context())

	s.dbStore.LoggedDB.Lock()
	defer s.dbStore.LoggedDB.Unlock()

	u, err := s.dbStore.LoggedDB.Get(username)
	if err != nil {
		httputil.WriteError(w, er.ErrNotLogged, http.StatusBadRequest)
		return
	}

	u.LastActivity = time.Now()
}
