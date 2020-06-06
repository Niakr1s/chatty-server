package server

import (
	"encoding/json"
	"net/http"

	"github.com/niakr1s/chatty-server/app/internal/httputil"
	"github.com/niakr1s/chatty-server/app/models"
)

// CreateChat ...
func (s *Server) CreateChat(w http.ResponseWriter, r *http.Request) {
	req := models.Chat{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		httputil.WriteError(w, err, http.StatusBadRequest)
		return
	}

	if err := s.dbStore.ChatDB.Add(req.ChatName); err != nil {
		httputil.WriteError(w, err, http.StatusBadRequest)
		return
	}
}
