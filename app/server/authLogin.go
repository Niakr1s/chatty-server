package server

import (
	"encoding/json"
	"net/http"

	"github.com/niakr1s/chatty-server/app/constants"
	"github.com/niakr1s/chatty-server/app/internal/httputil"
	"github.com/niakr1s/chatty-server/app/internal/sess"
	"github.com/niakr1s/chatty-server/app/models"
)

// AuthLogin should be used always after AuthOnly middleware
func (s *Server) AuthLogin(w http.ResponseWriter, r *http.Request) {

	session, err := sess.GetSessionFromStore(s.cookieStore, r)
	if err != nil {
		httputil.WriteSessionError(w)
		return
	}

	username := sess.GetUserNameFromCtx(r.Context())

	s.dbStore.LoggedDB.Logout(username) // we are authorized, do forced login
	loggedu, err := s.dbStore.LoggedDB.Login(username)
	if err != nil {
		httputil.WriteError(w, err, http.StatusForbidden)
		return
	}

	if storedU, err := s.dbStore.UserDB.Get(loggedu.UserName); err == nil {
		loggedu.UserStatus = storedU.UserStatus
		s.dbStore.LoggedDB.Update(loggedu)
	}

	session.Values[constants.SessionUserName] = loggedu.UserName
	session.Values[constants.SessionLoginToken] = loggedu.LoginToken

	if err := session.Save(r, w); err != nil {
		httputil.WriteSessionError(w)
		return
	}

	userToReturn := models.NewUserWithStatus(loggedu.User, loggedu.UserStatus)

	if err := json.NewEncoder(w).Encode(userToReturn); err != nil {
		httputil.WriteError(w, err, http.StatusInternalServerError)
		return
	}
}
