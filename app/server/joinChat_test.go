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

func TestServer_JoinChat(t *testing.T) {
	const chatname = "chat"
	const username = "user"

	testCases := []struct {
		name       string
		username   string
		req        string
		okExpected bool
	}{
		{
			"other user",
			"other user",
			fmt.Sprintf(`{"chat":"%s"}`, chatname),
			true,
		},
		{
			"same user",
			"user",
			fmt.Sprintf(`{"chat":"%s"}`, chatname),
			false,
		},
		{
			"other user non-existent chat",
			"other user",
			fmt.Sprintf(`{"chat":"%s"}`, chatname+"1"),
			false,
		},
		{
			"other user empty chat",
			"other user",
			fmt.Sprintf(`{"chat":""}`),
			false,
		},
		{
			"same user non-existent chat",
			username,
			fmt.Sprintf(`{"chat":"%s"}`, chatname+"1"),
			false,
		},
		{
			"other user bad request",
			"other user",
			fmt.Sprintf(`{"chatasf":"%s"}`, chatname),
			false,
		},
		{
			"other user bad request",
			"other user",
			"",
			false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			s := newMockServer()
			chat, _ := s.dbStore.ChatDB.Add(chatname)
			chat.AddUser(username)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.req))
			r = r.WithContext(sess.SetUserNameIntoCtx(r.Context(), tt.username))

			s.JoinChat(w, r)

			assert.Equal(t, tt.okExpected, w.Code == http.StatusOK)
		})
	}

}
