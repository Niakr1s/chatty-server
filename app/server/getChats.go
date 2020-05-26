package server

import (
	"encoding/json"
	"net/http"

	"github.com/niakr1s/chatty-server/app/internal/httputil"
)

// GetChats ...
func (s *Server) GetChats(w http.ResponseWriter, r *http.Request) {
	s.dbStore.ChatDB.Lock()
	defer s.dbStore.ChatDB.Unlock()

	chats := s.dbStore.ChatDB.GetChats()

	if err := json.NewEncoder(w).Encode(chats); err != nil {
		httputil.WriteError(w, err, http.StatusInternalServerError)
		return
	}
}
