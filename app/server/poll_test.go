package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/niakr1s/chatty-server/app/config"
	"github.com/stretchr/testify/assert"
)

func TestServer_Poll(t *testing.T) {
	s := NewMemoryServer()

	const chatname = "chat"
	const username = "user"

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
	r = r.WithContext(context.WithValue(r.Context(), config.CtxUserNameKey, username))

	done := make(chan struct{})
	go func() {
		s.Poll(w, r)
		s.Poll(w, r)
		done <- struct{}{}
	}()
	s.Store.ChatDB.Lock()

	chat, _ := s.Store.ChatDB.Add(chatname)
	chat.AddUser(username)
	chat.AddUser("another user")

	s.Store.ChatDB.Unlock()

	<-done

	assert.Equal(t, http.StatusOK, w.Code)
}
