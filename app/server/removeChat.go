package server

import (
	"encoding/json"
	"net/http"

	"github.com/niakr1s/chatty-server/app/internal/httputil"
	"github.com/niakr1s/chatty-server/app/models"
)

// RemoveChat ...
func (s *Server) RemoveChat(w http.ResponseWriter, r *http.Request) {
	req := models.Chat{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		httputil.WriteError(w, err, http.StatusBadRequest)
		return
	}

	s.dbStore.ChatDB.Lock()
	defer s.dbStore.ChatDB.Unlock()

	if err := s.dbStore.ChatDB.Remove(req.ChatName); err != nil {
		httputil.WriteError(w, err, http.StatusBadRequest)
		return
	}
}
