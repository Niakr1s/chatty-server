package server

import (
	"encoding/json"
	"net/http"

	"github.com/niakr1s/chatty-server/app/db"
	"github.com/niakr1s/chatty-server/app/er"
	"github.com/niakr1s/chatty-server/app/internal/httputil"
	"github.com/niakr1s/chatty-server/app/internal/sess"
	"github.com/niakr1s/chatty-server/app/internal/validator"
	"github.com/niakr1s/chatty-server/app/models"
)

// PostMessage ...
func (s *Server) PostMessage(w http.ResponseWriter, r *http.Request) {
	mess := &models.Message{}
	err := json.NewDecoder(r.Body).Decode(mess)
	if err != nil {
		httputil.WriteError(w, err, http.StatusBadRequest)
		return
	}

	mess.UserName = sess.GetUserNameFromCtx(r.Context())

	if err := validator.Validate.Struct(mess); err != nil {
		httputil.WriteError(w, err, http.StatusBadRequest)
		return
	}

	s.dbStore.ChatDB.Lock()
	defer s.dbStore.ChatDB.Unlock()

	chat, err := s.dbStore.ChatDB.Get(mess.ChatName)
	if err != nil {
		httputil.WriteError(w, er.ErrNoSuchChat, http.StatusBadRequest)
		return
	}

	chat.Lock()
	defer chat.Unlock()

	if !chat.IsInChat(mess.UserName) {
		httputil.WriteError(w, er.ErrNotInChat, http.StatusBadRequest)
		return
	}

	s.dbStore.MessageDB.Lock()
	defer s.dbStore.MessageDB.Unlock()

	mess.Verified = db.IsUserVerified(s.dbStore.UserDB, mess.UserName)

	err = s.dbStore.MessageDB.Post(mess)
	if err != nil {
		httputil.WriteError(w, err, http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(mess); err != nil {
		httputil.WriteError(w, err, http.StatusInternalServerError)
		return
	}
}
