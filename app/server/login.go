package server

import (
	"encoding/json"
	"net/http"

	"github.com/niakr1s/chatty-server/app/constants"
	"github.com/niakr1s/chatty-server/app/er"
	"github.com/niakr1s/chatty-server/app/internal/httputil"
	"github.com/niakr1s/chatty-server/app/internal/sess"
	"github.com/niakr1s/chatty-server/app/internal/validator"
	"github.com/niakr1s/chatty-server/app/models"
)

// Login - without password!
func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	u := models.User{}
	err := json.NewDecoder(r.Body).Decode(&u)

	if err != nil {
		httputil.WriteError(w, er.ErrCannotParseData, http.StatusBadRequest)
		return
	}

	if err := validator.Validate.Struct(u); err != nil {
		httputil.WriteError(w, err, http.StatusBadRequest)
		return
	}

	session, err := sess.GetSessionFromStore(s.cookieStore, r)
	if err != nil {
		httputil.WriteSessionError(w)
		return
	}

	if fullU, err := s.dbStore.UserDB.Get(u.UserName); err == nil || fullU.Verified {
		httputil.WriteError(w, er.ErrUnathorized, http.StatusConflict)
		return
	}

	s.dbStore.LoggedDB.Lock()
	defer s.dbStore.LoggedDB.Unlock()

	loggedU, err := s.dbStore.LoggedDB.Login(u.UserName)
	if err != nil {
		httputil.WriteError(w, err, http.StatusConflict)
		return
	}

	session.Values[constants.SessionUserName] = loggedU.UserName
	session.Values[constants.SessionLoginToken] = loggedU.LoginToken

	if err := session.Save(r, w); err != nil {
		httputil.WriteSessionError(w)
		return
	}

	if err := json.NewEncoder(w).Encode(loggedU.User); err != nil {
		httputil.WriteError(w, err, http.StatusInternalServerError)
		return
	}
}
