package server

import (
	"encoding/json"
	"net/http"

	"github.com/niakr1s/chatty-server/app/config"
	"github.com/niakr1s/chatty-server/app/er"
	"github.com/niakr1s/chatty-server/app/internal/httputil"
	"github.com/niakr1s/chatty-server/app/internal/sess"
	"github.com/niakr1s/chatty-server/app/models"
)

// GetLastMessages ...
func (s *Server) GetLastMessages(w http.ResponseWriter, r *http.Request) {
	username := sess.GetUserNameFromCtx(r.Context())

	input := models.Chat{}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		httputil.WriteError(w, err, http.StatusBadRequest)
		return
	}

	if !s.dbStore.ChatDB.IsInChat(input.ChatName, username) {
		httputil.WriteError(w, er.ErrNotInChat, http.StatusBadRequest)
		return
	}

	res, _ := s.dbStore.MessageDB.GetLastNMessages(input.ChatName, config.C.LastMessages)

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		httputil.WriteError(w, err, http.StatusInternalServerError)
		return
	}
}
