package server

import (
	"encoding/json"
	"net/http"

	"github.com/niakr1s/chatty-server/app/er"
	"github.com/niakr1s/chatty-server/app/models"
)

// Register ...
func (s *Server) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	u := models.User{}
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		s.writeError(w, er.ErrCannotParseData, http.StatusBadRequest)
		return
	}

	if err := u.GeneratePasswordHash(); err != nil {
		s.writeError(w, err, http.StatusUnprocessableEntity)
		return
	}

	if err := u.ValidateBeforeStoring(); err != nil {
		s.writeError(w, err, http.StatusBadRequest)
		return
	}

	if err := s.store.UserDB.Store(&u); err != nil {
		s.writeError(w, err, http.StatusConflict)
		return
	}

	res := struct {
		ID   interface{} `json:"id"`
		Name interface{} `json:"name"`
	}{u.ID, u.Name}

	json.NewEncoder(w).Encode(res)
}
