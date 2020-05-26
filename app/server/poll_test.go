package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/niakr1s/chatty-server/app/constants"
	"github.com/stretchr/testify/assert"
)

func TestServer_Poll(t *testing.T) {
	s := NewMemoryServer()

	const chatname = "chat"
	const username = "user"

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
	r = r.WithContext(context.WithValue(r.Context(), constants.CtxUserNameKey, username))

	done := make(chan struct{})
	go func() {
		s.Poll(w, r)
		s.Poll(w, r)
		done <- struct{}{}
	}()
	s.dbStore.ChatDB.Lock()

	chat, _ := s.dbStore.ChatDB.Add(chatname)
	chat.AddUser(username)
	chat.AddUser("another user")

	s.dbStore.ChatDB.Unlock()

	<-done

	assert.Equal(t, http.StatusOK, w.Code)
}
