package server

import (
	"encoding/json"
	"net/http"

	"github.com/niakr1s/chatty-server/app/internal/httputil"
	"github.com/niakr1s/chatty-server/app/internal/sess"
)

// GetChats ...
func (s *Server) GetChats(w http.ResponseWriter, r *http.Request) {
	username := sess.GetUserNameFromCtx(r.Context())

	s.dbStore.ChatDB.Lock()
	defer s.dbStore.ChatDB.Unlock()

	chats := s.dbStore.ChatDB.GetChats()

	type result struct {
		Chatname string `json:"name"`
		Joined   bool   `json:"joined"`
	}

	res := make([]result, 0, len(chats))

	for _, c := range chats {
		res = append(res, result{c.ChatName(), c.IsInChat(username)})
	}

	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		httputil.WriteError(w, err, http.StatusInternalServerError)
		return
	}
}
