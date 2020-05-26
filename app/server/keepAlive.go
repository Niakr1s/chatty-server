package server

import (
	"net/http"
	"time"

	"github.com/niakr1s/chatty-server/app/constants"
	"github.com/niakr1s/chatty-server/app/er"
	"github.com/niakr1s/chatty-server/app/internal/httputil"
)

// KeepAlive ...
func (s *Server) KeepAlive(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(constants.CtxUserNameKey).(string)

	s.dbStore.LoggedDB.Lock()
	defer s.dbStore.LoggedDB.Unlock()

	u, err := s.dbStore.LoggedDB.Get(username)
	if err != nil {
		httputil.WriteError(w, er.ErrNotLogged, http.StatusBadRequest)
		return
	}

	u.LastActivity = time.Now()
}
