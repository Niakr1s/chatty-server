package server

import (
	"encoding/json"
	"net/http"

	"github.com/niakr1s/chatty-server/app/internal/httputil"
	"github.com/niakr1s/chatty-server/app/internal/sess"
)

// Poll ...
func (s *Server) Poll(w http.ResponseWriter, r *http.Request) {
	username := sess.GetUserNameFromCtx(r.Context())

	event := <-s.pool.GetUserChan(username)

	if err := json.NewEncoder(w).Encode(event); err != nil {
		httputil.WriteError(w, err, http.StatusInternalServerError)
		return
	}
}
