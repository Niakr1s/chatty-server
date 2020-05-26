package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/niakr1s/chatty-server/app/er"
	"github.com/niakr1s/chatty-server/app/internal/httputil"
)

// VerifyEmail .../{username}/{activationToken}
func (s *Server) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	username, ok := vars["username"]
	if !ok {
		httputil.WriteError(w, er.ErrNoUsername, http.StatusBadRequest)
		return
	}

	activationToken, ok := vars["activationToken"]
	if !ok {
		httputil.WriteError(w, er.ErrNoActivationToken, http.StatusBadRequest)
		return
	}

	u, err := s.dbStore.UserDB.Get(username)
	if err != nil {
		httputil.WriteError(w, er.ErrUserNotFound, http.StatusBadRequest)
		return
	}

	if u.Email.ActivationToken != activationToken {
		httputil.WriteError(w, er.ErrBadActivationToken, http.StatusBadRequest)
		return
	}

	u.Activated = true

	if err = s.dbStore.UserDB.Update(&u); err != nil {
		httputil.WriteError(w, er.ErrCannotUpdateUser, http.StatusInternalServerError)
		return
	}
}
