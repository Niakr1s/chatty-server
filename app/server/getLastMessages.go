package server

import (
	"encoding/json"
	"net/http"

	"github.com/niakr1s/chatty-server/app/config"
	"github.com/niakr1s/chatty-server/app/er"
	"github.com/niakr1s/chatty-server/app/internal/httputil"
	"github.com/niakr1s/chatty-server/app/internal/sess"
)

// GetLastMessages ...
func (s *Server) GetLastMessages(w http.ResponseWriter, r *http.Request) {
	username := sess.GetUserNameFromCtx(r.Context())

	input := struct {
		Chatname string `json:"chatname"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		httputil.WriteError(w, err, http.StatusBadRequest)
		return
	}

	s.dbStore.ChatDB.Lock()
	defer s.dbStore.ChatDB.Unlock()

	c, err := s.dbStore.ChatDB.Get(input.Chatname)
	if err != nil {
		httputil.WriteError(w, err, http.StatusBadRequest)
		return
	}

	c.Lock()
	defer c.Unlock()

	if !c.IsInChat(username) {
		httputil.WriteError(w, er.ErrNotInChat, http.StatusBadRequest)
		return
	}

	s.dbStore.MessageDB.Lock()
	defer s.dbStore.MessageDB.Unlock()

	res, _ := s.dbStore.MessageDB.GetLastNMessages(c.ChatName(), config.C.LastMessages)

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		httputil.WriteError(w, err, http.StatusInternalServerError)
		return
	}
}
