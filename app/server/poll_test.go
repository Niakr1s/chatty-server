package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/niakr1s/chatty-server/app/internal/sess"
	"github.com/stretchr/testify/assert"
)

func TestServer_Poll(t *testing.T) {
	s := NewMemoryServer()

	const chatname = "chat"
	const username = "user"

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
	r = r.WithContext(sess.SetUserNameIntoCtx(r.Context(), username))

	done := make(chan struct{})
	go func() {
		<-time.After(time.Millisecond * 10)
		s.Poll(w, r)
		s.Poll(w, r)
		s.Poll(w, r)
		done <- struct{}{}
	}()

	s.dbStore.LoggedDB.Login(username)
	s.dbStore.ChatDB.Lock()

	chat, _ := s.dbStore.ChatDB.Add(chatname)
	chat.AddUser(username)
	chat.AddUser("another user")

	s.dbStore.ChatDB.Unlock()

	<-done

	assert.Equal(t, http.StatusOK, w.Code)
}
