package server

import (
	"encoding/json"
	"net/http"

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

	if !s.dbStore.ChatDB.IsInChat(mess.ChatName, mess.UserName) {
		httputil.WriteError(w, er.ErrNotInChat, http.StatusBadRequest)
		return
	}

	if storedU, err := s.dbStore.UserDB.Get(mess.UserName); err == nil {
		mess.UserStatus = storedU.UserStatus
	}

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
