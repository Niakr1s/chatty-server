package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/niakr1s/chatty-server/app/models"
	"github.com/stretchr/testify/assert"
)

func TestServer_RemoveChat(t *testing.T) {
	const chatname = "chat"

	tests := []struct {
		name                string
		chat                string
		okExpected          bool
		expectedChatsLength int
	}{
		{"same chat", chatname, true, 0},
		{"non-existent chat", "other chat", false, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newMockServer()

			s.dbStore.ChatDB.Add(chatname)

			req := models.Chat{ChatName: tt.chat}
			b, _ := json.Marshal(req)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))

			s.RemoveChat(w, r)

			assert.Equal(t, tt.okExpected, w.Code == http.StatusOK)
			assert.Len(t, s.dbStore.ChatDB.GetChats(), tt.expectedChatsLength)
		})
	}

}
