package server

import (
	"encoding/json"
	"net/http"

	"github.com/niakr1s/chatty-server/app/er"
	"github.com/niakr1s/chatty-server/app/internal/httputil"
	"github.com/niakr1s/chatty-server/app/internal/sess"
)

// LeaveChat ...
func (s *Server) LeaveChat(w http.ResponseWriter, r *http.Request) {
	username := sess.GetUserNameFromCtx(r.Context())

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

	if err := chat.RemoveUser(username); err != nil {
		httputil.WriteError(w, err, http.StatusBadRequest)
		return
	}
}
