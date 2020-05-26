package server

import (
	"encoding/json"
	"net/http"

	"github.com/niakr1s/chatty-server/app/constants"
	"github.com/niakr1s/chatty-server/app/er"
	"github.com/niakr1s/chatty-server/app/internal/httputil"
)

// JoinChat ...
func (s *Server) JoinChat(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(constants.CtxUserNameKey).(string)

	req := struct {
		Chatname string `json:"chatname"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		httputil.WriteError(w, er.ErrCannotParseData, http.StatusBadRequest)
		return
	}

	s.dbStore.ChatDB.Lock()
	defer s.dbStore.ChatDB.Unlock()

	chat, err := s.dbStore.ChatDB.Get(req.Chatname)
	if err != nil {
		httputil.WriteError(w, err, http.StatusBadRequest)
		return
	}

	chat.Lock()
	defer chat.Unlock()

	if chat.IsInChat(username) {
		httputil.WriteError(w, er.ErrAlreadyInChat, http.StatusBadRequest)
		return
	}

	if err := chat.AddUser(username); err != nil {
		httputil.WriteError(w, err, http.StatusBadRequest)
		return
	}
}
