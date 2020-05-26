package server

import (
	"encoding/json"
	"net/http"

	"github.com/niakr1s/chatty-server/app/config"
	"github.com/niakr1s/chatty-server/app/er"
	"github.com/niakr1s/chatty-server/app/models"
	"github.com/niakr1s/chatty-server/app/server/httputil"
	"github.com/niakr1s/chatty-server/app/server/sess"
)

// Authorize ...
func (s *Server) Authorize(w http.ResponseWriter, r *http.Request) {
	u := models.FullUser{}
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		httputil.WriteError(w, er.ErrCannotParseData, http.StatusBadRequest)
		return
	}

	storedU, err := s.store.UserDB.Get(u.Name)
	if err != nil {
		httputil.WriteError(w, er.ErrUserNotFound, http.StatusBadRequest)
		return
	}

	if !storedU.Email.Activated {
		httputil.WriteError(w, er.ErrUnverifiedEmail, http.StatusForbidden)
		return
	}

	if err := storedU.CheckPassword(u.Password); err != nil {
		httputil.WriteError(w, er.ErrHashMismatch, http.StatusUnauthorized)
		return
	}

	session, err := sess.GetSessionFromStore(s.cookieStore, r)
	if err != nil {
		httputil.WriteSessionError(w)
		return
	}

	s.store.LoggedDB.Lock()
	defer s.store.LoggedDB.Unlock()

	err = s.store.LoggedDB.Logout(storedU.Name) // force logout, don't check err

	loggedU, err := s.store.LoggedDB.Login(storedU.Name)
	if err != nil {
		httputil.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	session.Values[config.SessionUserName] = loggedU.Name
	session.Values[config.SessionLoginToken] = loggedU.LoginToken

	if err := session.Save(r, w); err != nil {
		httputil.WriteSessionError(w)
		return
	}
}
