package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/niakr1s/chatty-server/app/internal/sess"
	"github.com/niakr1s/chatty-server/app/models"
	"github.com/stretchr/testify/assert"
)

func TestServer_Authorize(t *testing.T) {
	s := newMockServer()

	storedUser := models.NewFullUser("user", "user@example.org", "password")
	storedUser.GeneratePasswordHash()
	storedUser.Verified = true

	s.dbStore.UserDB.Store(storedUser)

	unverifiedUser := models.NewFullUser("unverifiedUser", "user@example.org", "password")
	s.dbStore.UserDB.Store(unverifiedUser)

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
		{"user with unverified email", "unverifiedUser", "password", false},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			u := models.NewFullUser(tt.username, "user@example.org", tt.password)
			b, _ := json.Marshal(&u)
			r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))
			w := httptest.NewRecorder()

			s.Authorize(w, r)

			assert.Equal(t, tt.okExpected, w.Code == http.StatusOK)

			r.Header.Set("Cookie", w.HeaderMap.Get("Set-Cookie"))

			session, _ := sess.GetSessionFromStore(s.cookieStore, r)
			username, _ := sess.GetUserName(session)

			if tt.okExpected {
				assert.Equal(t, tt.username, username)
			}
		})
	}
}
