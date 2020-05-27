package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/niakr1s/chatty-server/app/internal/sess"
	"github.com/stretchr/testify/assert"
)

func TestServer_GetLastMessages(t *testing.T) {
	const username = "user"
	const chat = "chat"

	tests := []struct {
		name       string
		user       string
		chat       string
		okExpected bool
	}{
		{"valid", username, chat, true},
		{"unlogged chat", username, "unlogged chat", false},
		{"unlogged user", "unlogged user", chat, false},
		{"unlogged user unlogged chat", "unlogged user", "unlogged chat", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewMemoryServer()
			chat, _ := s.dbStore.ChatDB.Add(chat)
			chat.AddUser(username)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(fmt.Sprintf(`{"chatname":"%s"}`, tt.chat)))

			r = r.WithContext(sess.SetUserNameIntoCtx(r.Context(), tt.user))

			s.GetLastMessages(w, r)

			assert.Equal(t, tt.okExpected, w.Code == http.StatusOK)
		})
	}
}
