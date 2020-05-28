package server

import (
	"encoding/json"
	"net/http"

	"github.com/niakr1s/chatty-server/app/constants"
	"github.com/niakr1s/chatty-server/app/internal/httputil"
	"github.com/niakr1s/chatty-server/app/internal/sess"
)

// AuthLogin should be used always after AuthOnly middleware
func (s *Server) AuthLogin(w http.ResponseWriter, r *http.Request) {

	session, err := sess.GetSessionFromStore(s.cookieStore, r)
	if err != nil {
		httputil.WriteSessionError(w)
		return
	}

	username := sess.GetUserNameFromCtx(r.Context())

	s.dbStore.LoggedDB.Lock()
	defer s.dbStore.LoggedDB.Unlock()

	s.dbStore.LoggedDB.Logout(username) // we are authorized, do forced login
	u, err := s.dbStore.LoggedDB.Login(username)

	if err != nil {
		httputil.WriteError(w, err, http.StatusForbidden)
		return
	}

	session.Values[constants.SessionLoginToken] = u.LoginToken
	session.Save(r, w)

	if err := json.NewEncoder(w).Encode(u.User); err != nil {
		httputil.WriteError(w, err, http.StatusInternalServerError)
		return
	}
}
