package server

import (
	"encoding/json"
	"net/http"

	"github.com/niakr1s/chatty-server/app/er"
	"github.com/niakr1s/chatty-server/app/internal/httputil"
	"github.com/niakr1s/chatty-server/app/internal/sess"
	"github.com/niakr1s/chatty-server/app/models"
)

// GetUsers ...
func (s *Server) GetUsers(w http.ResponseWriter, r *http.Request) {
	username := sess.GetUserNameFromCtx(r.Context())

	req := models.Chat{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		httputil.WriteError(w, err, http.StatusBadRequest)
		return
	}

	if !s.dbStore.ChatDB.IsInChat(req.ChatName, username) {
		httputil.WriteError(w, er.ErrNotInChat, http.StatusBadRequest)
		return
	}

	users := s.dbStore.ChatDB.GetUsers(req.ChatName)
	loggedUsers := make([]*models.LoggedUser, 0, len(users))

	for _, u := range users {
		lu, err := s.dbStore.LoggedDB.Get(u)
		if err != nil {
			continue
		}
		loggedUsers = append(loggedUsers, lu)
	}

	if err := json.NewEncoder(w).Encode(loggedUsers); err != nil {
		httputil.WriteError(w, err, http.StatusInternalServerError)
		return
	}
}
