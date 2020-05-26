package server

import (
	"encoding/json"
	"net/http"

	"github.com/niakr1s/chatty-server/app/constants"
	"github.com/niakr1s/chatty-server/app/er"
	"github.com/niakr1s/chatty-server/app/internal/httputil"
	"github.com/niakr1s/chatty-server/app/internal/sess"
	"github.com/niakr1s/chatty-server/app/models"
)

// Authorize ...
func (s *Server) Authorize(w http.ResponseWriter, r *http.Request) {
	u := models.FullUser{}
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		httputil.WriteError(w, er.ErrCannotParseData, http.StatusBadRequest)
		return
	}

	storedU, err := s.dbStore.UserDB.Get(u.Name)
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

	s.dbStore.LoggedDB.Lock()
	defer s.dbStore.LoggedDB.Unlock()

	err = s.dbStore.LoggedDB.Logout(storedU.Name) // force logout, don't check err

	loggedU, err := s.dbStore.LoggedDB.Login(storedU.Name)
	if err != nil {
		httputil.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	session.Values[constants.SessionUserName] = loggedU.Name
	session.Values[constants.SessionLoginToken] = loggedU.LoginToken

	if err := session.Save(r, w); err != nil {
		httputil.WriteSessionError(w)
		return
	}
}
