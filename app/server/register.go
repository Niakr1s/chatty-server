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

// Register ...
func (s *Server) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	u := models.User{}
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		httputil.WriteError(w, er.ErrCannotParseData, http.StatusBadRequest)
		return
	}

	if err := u.GeneratePasswordHash(); err != nil {
		httputil.WriteError(w, err, http.StatusUnprocessableEntity)
		return
	}

	if err := u.ValidateBeforeStoring(); err != nil {
		httputil.WriteError(w, err, http.StatusBadRequest)
		return
	}

	if err := s.store.UserDB.Store(&u); err != nil {
		httputil.WriteError(w, err, http.StatusConflict)
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
