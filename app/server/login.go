package server

import (
	"encoding/json"
	"net/http"

	"github.com/niakr1s/chatty-server/app/config"
	"github.com/niakr1s/chatty-server/app/er"
	"github.com/niakr1s/chatty-server/app/models"
	"github.com/niakr1s/chatty-server/app/server/httputil"
	"github.com/niakr1s/chatty-server/app/server/sess"
	"github.com/niakr1s/chatty-server/app/validator"
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

	s.store.LoggedDB.Lock()
	defer s.store.LoggedDB.Unlock()

	loggedU, err := s.store.LoggedDB.Login(u.Name)
	if err != nil {
		httputil.WriteError(w, err, http.StatusConflict)
		return
	}

	session.Values[config.SessionUserName] = loggedU.Name
	session.Values[config.SessionLoginToken] = loggedU.LoginToken

	if err := session.Save(r, w); err != nil {
		httputil.WriteSessionError(w)
		return
	}

	json.NewEncoder(w).Encode(loggedU.User)
}
