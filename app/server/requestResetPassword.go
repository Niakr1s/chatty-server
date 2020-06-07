package server

import (
	"encoding/json"
	"net/http"

	"github.com/niakr1s/chatty-server/app/er"
	"github.com/niakr1s/chatty-server/app/internal/httputil"
	"github.com/niakr1s/chatty-server/app/models"
)

// RequestResetPassword ...
func (s *Server) RequestResetPassword(w http.ResponseWriter, r *http.Request) {
	uwe := models.User{}
	if err := json.NewDecoder(r.Body).Decode(&uwe); err != nil {
		httputil.WriteError(w, err, http.StatusBadRequest)
		return
	}

	if uwe.UserName == "" {
		httputil.WriteError(w, er.ErrUserNameIsEmpty, http.StatusBadRequest)
		return
	}

	fullU, err := s.dbStore.UserDB.Get(uwe.UserName)
	if err != nil {
		httputil.WriteError(w, er.ErrUserNotFound, http.StatusInternalServerError)
		return
	}

	fullU.Pass.GeneratePasswordResetToken()
	if err := s.dbStore.UserDB.Update(fullU); err != nil {
		httputil.WriteError(w, er.ErrCannotUpdateUser, http.StatusInternalServerError)
		return
	}

	if err := s.mailer.SendResetPasswordEmail(fullU.Email.Address, fullU.UserName, fullU.Pass.PasswordResetToken); err != nil {
		httputil.WriteError(w, er.ErrSendEmail, http.StatusInternalServerError)
		return
	}
}
