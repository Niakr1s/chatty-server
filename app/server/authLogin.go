package server

import (
	"encoding/json"
	"net/http"

	"github.com/niakr1s/chatty-server/app/config"
	"github.com/niakr1s/chatty-server/app/er"
	"github.com/niakr1s/chatty-server/app/server/httputil"
	"github.com/niakr1s/chatty-server/app/server/sess"
)

// AuthLogin should be used always after AuthOnly middleware
func (s *Server) AuthLogin(w http.ResponseWriter, r *http.Request) {

	session, err := sess.GetSessionFromStore(s.cookieStore, r)
	if err != nil {
		httputil.WriteSessionError(w)
		return
	}

	username, ok := r.Context().Value(config.CtxUserNameKey).(string)
	if !ok {
		httputil.WriteError(w, er.ErrUserNameIsEmpty, http.StatusForbidden)
		return
	}

	s.Store.LoggedDB.Lock()
	defer s.Store.LoggedDB.Unlock()

	s.Store.LoggedDB.Logout(username) // we are authorized, do forced login
	u, err := s.Store.LoggedDB.Login(username)

	if err != nil {
		httputil.WriteError(w, err, http.StatusForbidden)
		return
	}

	session.Values[config.SessionLoginToken] = u.LoginToken
	session.Save(r, w)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u.User)
}
