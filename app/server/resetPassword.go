package server

import (
	"encoding/json"
	"net/http"

	"github.com/niakr1s/chatty-server/app/er"
	"github.com/niakr1s/chatty-server/app/internal/httputil"
	"github.com/niakr1s/chatty-server/app/internal/validator"
	"github.com/niakr1s/chatty-server/app/models"
)

// ResetPassword ...
func (s *Server) ResetPassword(w http.ResponseWriter, r *http.Request) {
	type userWithPass struct {
		models.User
		models.Pass
	}
	req := userWithPass{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, err, http.StatusBadRequest)
		return
	}

	if req.UserName == "" || req.PasswordResetToken == "" || req.Password == "" {
		httputil.WriteError(w, er.ErrEmptyFields, http.StatusBadRequest)
		return
	}

	fullU, err := s.dbStore.UserDB.Get(req.UserName)
	if err != nil {
		httputil.WriteError(w, er.ErrUserNotFound, http.StatusInternalServerError)
		return
	}

	if fullU.PasswordResetToken != req.PasswordResetToken {
		httputil.WriteError(w, er.ErrBadResetPasswordToken, http.StatusBadRequest)
		return
	}

	if err := req.Pass.GeneratePasswordHash(); err != nil {
		httputil.WriteError(w, err, http.StatusUnprocessableEntity)
		return
	}

	fullU.Pass = req.Pass

	if err := validator.Validate.Struct(fullU); err != nil {
		httputil.WriteError(w, err, http.StatusInternalServerError)
		return
	}

	fullU.EraseResetToken()
	fullU.ErasePassword()

	if err := s.dbStore.UserDB.Update(fullU); err != nil {
		httputil.WriteError(w, er.ErrCannotUpdateUser, http.StatusInternalServerError)
		return
	}
}
