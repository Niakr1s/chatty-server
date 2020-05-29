package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/niakr1s/chatty-server/app/db"
	"github.com/niakr1s/chatty-server/app/internal/sess"
	"github.com/niakr1s/chatty-server/app/models"
	"github.com/stretchr/testify/assert"
)

func TestServer_GetChats(t *testing.T) {
	const username = "user"

	tests := []struct {
		name       string
		username   string
		okExpected bool
	}{
		{"valid", username, true},
		{"non logged user", "non logged", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewMemoryServer()
			c, _ := s.dbStore.ChatDB.Add("chat")
			s.dbStore.LoggedDB.Login(username)
			c.AddUser(username)
			s.dbStore.MessageDB.Post(models.NewMessage(username, "text", "chat"))

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))

			r = r.WithContext(sess.SetUserNameIntoCtx(r.Context(), tt.username))

			s.GetChats(w, r)

			assert.Equal(t, tt.okExpected, w.Code == http.StatusOK)

			report := []db.ChatReport{}
			json.NewDecoder(w.Body).Decode(&report)

			assert.Equal(t, 1, len(report))
			assert.Equal(t, tt.username == username, report[0].Joined)
			if tt.username == username {
				assert.NotEmpty(t, report[0].Users)
				assert.NotEmpty(t, report[0].Messages)
			} else {
				assert.Empty(t, report[0].Users)
				assert.Empty(t, report[0].Messages)
			}
		})
	}
}
