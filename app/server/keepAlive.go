package server

import (
	"net/http"
	"time"

	"github.com/niakr1s/chatty-server/app/config"
	"github.com/niakr1s/chatty-server/app/er"
	"github.com/niakr1s/chatty-server/app/server/httputil"
)

// KeepAlive ...
func (s *Server) KeepAlive(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(config.CtxUserNameKey).(string)

	s.Store.LoggedDB.Lock()
	defer s.Store.LoggedDB.Unlock()

	u, err := s.Store.LoggedDB.Get(username)
	if err != nil {
		httputil.WriteError(w, er.ErrNotLogged, http.StatusBadRequest)
		return
	}

	u.LastActivity = time.Now()
}