package server

import (
	"encoding/json"
	"net/http"

	"github.com/niakr1s/chatty-server/app/er"
	"github.com/niakr1s/chatty-server/app/internal/httputil"
	"github.com/niakr1s/chatty-server/app/internal/sess"
	"github.com/niakr1s/chatty-server/app/models"
)

// JoinChat ...
func (s *Server) JoinChat(w http.ResponseWriter, r *http.Request) {
	username := sess.GetUserNameFromCtx(r.Context())

	req := models.Chat{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		httputil.WriteError(w, er.ErrCannotParseData, http.StatusBadRequest)
		return
	}

	if s.dbStore.ChatDB.IsInChat(req.ChatName, username) {
		httputil.WriteError(w, er.ErrAlreadyInChat, http.StatusBadRequest)
		return
	}

	if err := s.dbStore.ChatDB.AddUser(req.ChatName, username); err != nil {
		httputil.WriteError(w, err, http.StatusBadRequest)
		return
	}

	report := s.dbStore.MakeChatReportForUser(username, req.ChatName)

	if err := json.NewEncoder(w).Encode(report); err != nil {
		httputil.WriteError(w, err, http.StatusInternalServerError)
		return
	}
}
