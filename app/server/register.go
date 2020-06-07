package server

import (
	"encoding/json"
	"net/http"

	"github.com/niakr1s/chatty-server/app/er"
	"github.com/niakr1s/chatty-server/app/internal/httputil"
	"github.com/niakr1s/chatty-server/app/models"
)

// Register ...
func (s *Server) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	u := models.FullUser{}
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		httputil.WriteError(w, er.ErrCannotParseData, http.StatusBadRequest)
		return
	}

	u.GenerateActivationToken()

	if err := u.GeneratePasswordHash(); err != nil {
		httputil.WriteError(w, err, http.StatusUnprocessableEntity)
		return
	}

	if err := u.ValidateBeforeStoring(); err != nil {
		httputil.WriteError(w, err, http.StatusBadRequest)
		return
	}

	u.ErasePassword()

	if err := s.mailer.SendActivationEmail(u.Email.Address, u.UserName, u.Email.ActivationToken); err != nil {
		httputil.WriteError(w, er.ErrSendEmail, http.StatusInternalServerError)
		return
	}

	if err := s.dbStore.UserDB.Store(u); err != nil {
		httputil.WriteError(w, err, http.StatusConflict)
		return
	}
}
