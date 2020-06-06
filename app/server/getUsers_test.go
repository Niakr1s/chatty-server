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

func TestServer_GetUsers(t *testing.T) {
	const (
		username = "user"
		chatname = "chat"
	)

	tests := []struct {
		name       string
		user       string
		chat       string
		okExpected bool
	}{
		{"valid user", username, chatname, true},
		{"unlogged user", "unlogged user", chatname, false},
		{"invalid chat", username, "invalid chat", false},
		{"unlogged user invalid chat", "unlogged user", "invalid chat", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newMockServer()
			s.dbStore.ChatDB.Add(chatname)
			s.dbStore.ChatDB.AddUser(chatname, username)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(fmt.Sprintf(`{"chat":"%s"}`, tt.chat)))
			r = r.WithContext(sess.SetUserNameIntoCtx(r.Context(), tt.user))
			s.GetUsers(w, r)

			assert.Equal(t, tt.okExpected, w.Code == http.StatusOK)
		})
	}
}
