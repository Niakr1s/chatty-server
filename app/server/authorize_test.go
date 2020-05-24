package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/niakr1s/chatty-server/app/models"
	"github.com/niakr1s/chatty-server/app/server/sess"
	"github.com/stretchr/testify/assert"
)

func TestServer_Authorize(t *testing.T) {
	s := NewMemoryServer()

	storedUser := &models.User{Name: "user", Password: "password"}
	storedUser.GeneratePasswordHash()

	s.store.UserDB.Store(storedUser)

	testCases := []struct {
		name       string
		username   string
		password   string
		okExpected bool
	}{
		{"same user", "user", "password", true},
		{"same user with wrong pass", "user", "wrongpassword", false},
		{"same user with empty pass", "user", "", false},
		{"other user", "user1", "password", false},
		{"other user with empty pass", "user1", "", false},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			u := models.User{Name: tt.username, Password: tt.password}
			b, _ := json.Marshal(&u)
			r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))
			w := httptest.NewRecorder()

			s.Authorize(w, r)

			assert.Equal(t, tt.okExpected, w.Code == http.StatusOK)

			session, _ := sess.GetSessionFromStore(s.cookieStore, r)
			username, _ := sess.GetUserName(session)
			isAuthorized := sess.IsAuthorized(session)

			assert.Equal(t, tt.okExpected, isAuthorized)

			if isAuthorized {
				assert.Equal(t, tt.username, username)
			}
		})
	}
}
