package server

import (
	"encoding/json"
	"net/http"
	"server2/app/models"
)

// Register ...
func (s *Server) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	u := models.User{}
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		s.writeError(w, ErrCannotParseData, http.StatusBadRequest)
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
