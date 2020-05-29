package server

import (
	"encoding/json"
	"net/http"

	"github.com/niakr1s/chatty-server/app/config"
	"github.com/niakr1s/chatty-server/app/internal/httputil"
	"github.com/niakr1s/chatty-server/app/internal/sess"
	"github.com/niakr1s/chatty-server/app/models"
)

// GetChats ...
func (s *Server) GetChats(w http.ResponseWriter, r *http.Request) {
	username := sess.GetUserNameFromCtx(r.Context())

	s.dbStore.ChatDB.Lock()
	defer s.dbStore.ChatDB.Unlock()

	chats := s.dbStore.ChatDB.GetChats()

	type result struct {
		models.Chat
		Joined   bool              `json:"joined"`
		Messages []*models.Message `json:"messages"`
	}

	res := make([]result, 0, len(chats))

	for _, c := range chats {
		c.Lock()
		isInChat := c.IsInChat(username)
		messages := make([]*models.Message, 0)
		if isInChat {
			gotMessages, err := s.dbStore.MessageDB.GetLastNMessages(c.ChatName(), config.C.LastMessages)
			if err == nil {
				messages = gotMessages
			}
		}
		res = append(res, result{models.Chat{ChatName: c.ChatName()}, c.IsInChat(username), messages})
		c.Unlock()
	}

	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		httputil.WriteError(w, err, http.StatusInternalServerError)
		return
	}
}
