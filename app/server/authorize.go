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

	if err := storedU.CheckPassword(u.Password); err != nil {
		httputil.WriteError(w, err, http.StatusUnauthorized)
		return
	}

	session, err := sess.GetSessionFromStore(s.cookieStore, r)
	if err != nil {
		httputil.WriteSessionError(w)
		return
	}

	session.Values[config.SessionAuthorized] = true
	session.Values[config.SessionUserName] = u.Name
	session.Save(r, w)
}
