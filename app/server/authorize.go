package server

import (
	"encoding/json"
	"net/http"

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

	storedU, err := s.dbStore.UserDB.Get(u.UserName)
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

	r = r.WithContext(sess.SetUserNameIntoCtx(r.Context(), storedU.UserName))

	s.AuthLogin(w, r)
}
