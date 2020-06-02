package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/niakr1s/chatty-server/app/events"
	"github.com/niakr1s/chatty-server/app/internal/httputil"
	"github.com/niakr1s/chatty-server/app/internal/sess"
)

// Poll ...
func (s *Server) Poll(w http.ResponseWriter, r *http.Request) {
	username := sess.GetUserNameFromCtx(r.Context())

	select {
	case event, ok := <-s.pool.GetUserChan(username):
		if !ok {
			return
		}
		ewt := events.NewEventWithType(event)

		if err := json.NewEncoder(w).Encode(ewt); err != nil {
			httputil.WriteError(w, err, http.StatusInternalServerError)
			return
		}
	case <-time.After(time.Second * 10):
	}
}
