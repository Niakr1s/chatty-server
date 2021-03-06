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

func TestServer_LeaveChat(t *testing.T) {
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
			false,
		},
		{
			"same user",
			"user",
			fmt.Sprintf(`{"chat":"%s"}`, chatname),
			true,
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
			fmt.Sprintf(`{"chatfsadf":"%s"}`, chatname),
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
			s.dbStore.ChatDB.Add(chatname)
			s.dbStore.ChatDB.AddUser(chatname, username)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.req))
			r = r.WithContext(sess.SetUserNameIntoCtx(r.Context(), tt.username))

			s.LeaveChat(w, r)

			assert.Equal(t, tt.okExpected, w.Code == http.StatusOK)
		})
	}

}
