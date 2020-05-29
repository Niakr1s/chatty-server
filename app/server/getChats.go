package server

import (
	"encoding/json"
	"net/http"

	"github.com/niakr1s/chatty-server/app/db"
	"github.com/niakr1s/chatty-server/app/internal/httputil"
	"github.com/niakr1s/chatty-server/app/internal/sess"
)

// GetChats ...
func (s *Server) GetChats(w http.ResponseWriter, r *http.Request) {
	username := sess.GetUserNameFromCtx(r.Context())

	s.dbStore.ChatDB.Lock()
	defer s.dbStore.ChatDB.Unlock()

	chats := s.dbStore.ChatDB.GetChats()

	res := make([]db.ChatReport, 0, len(chats))

	for _, c := range chats {
		if report, err := s.dbStore.MakeChatReportForUser(username, c.ChatName()); err == nil {
			res = append(res, report)
		}
	}

	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		httputil.WriteError(w, err, http.StatusInternalServerError)
		return
	}
}
