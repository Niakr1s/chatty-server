package server

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/niakr1s/chatty-server/app/constants"
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
			fmt.Sprintf(`{"chatname":"%s"}`, chatname),
			true,
		},
		{
			"same user",
			"user",
			fmt.Sprintf(`{"chatname":"%s"}`, chatname),
			false,
		},
		{
			"other user non-existent chat",
			"other user",
			fmt.Sprintf(`{"chatname":"%s"}`, chatname+"1"),
			false,
		},
		{
			"other user empty chat",
			"other user",
			fmt.Sprintf(`{"chatname":""}`),
			false,
		},
		{
			"same user non-existent chat",
			username,
			fmt.Sprintf(`{"chatname":"%s"}`, chatname+"1"),
			false,
		},
		{
			"other user bad request",
			"other user",
			fmt.Sprintf(`{"chat":"%s"}`, chatname),
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
			s := NewMemoryServer()
			chat, _ := s.dbStore.ChatDB.Add(chatname)
			chat.AddUser(username)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.req))
			r = r.WithContext(context.WithValue(r.Context(), constants.CtxUserNameKey, tt.username))

			s.JoinChat(w, r)

			assert.Equal(t, tt.okExpected, w.Code == http.StatusOK)
		})
	}

}
