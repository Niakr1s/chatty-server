package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/niakr1s/chatty-server/app/internal/sess"
	"github.com/stretchr/testify/assert"
)

func TestServer_PostMessage(t *testing.T) {
	const chatname = "chat"
	const username = "user"

	type mess struct {
		Chat string `json:"chat"`
		Text string `json:"text"`
	}

	tests := []struct {
		name       string
		user       string
		message    *mess
		okExpected bool
	}{
		{"logged user to logged chat", username, &mess{Text: "text", Chat: chatname}, true},
		{"logged user to unlogged chat", username, &mess{Text: "text", Chat: "other chat"}, false},
		{"unlogged user to logged chat", "other user", &mess{Text: "text", Chat: chatname}, false},
		{"unlogged user to unlogged chat", "other user", &mess{Text: "text", Chat: "other chat"}, false},
		{"logged user to logged chat with invalid text", username, &mess{Text: "", Chat: chatname}, false},
		{"empty user to logged chat", "", &mess{Text: "text", Chat: chatname}, false},
		{"logged user to empty chat", username, &mess{Text: "text", Chat: ""}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := json.Marshal(tt.message)
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))
			r = r.WithContext(sess.SetUserNameIntoCtx(r.Context(), tt.user))

			s := NewMemoryServer()
			chat, _ := s.dbStore.ChatDB.Add(chatname)
			chat.AddUser(username)

			s.PostMessage(w, r)

			assert.Equal(t, tt.okExpected, w.Code == http.StatusOK)
		})
	}
}
