package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/niakr1s/chatty-server/app/config"
	"github.com/niakr1s/chatty-server/app/server/sess"
	"github.com/stretchr/testify/assert"
)

func TestAuthOnly(t *testing.T) {
	testCases := []struct {
		name          string
		authorized    bool
		username      string
		shouldExecute bool
	}{
		{"authorized user", true, "username", true},
		{"unauthorized user", false, "username", false},
		{"unauthorized user with empty name", false, "", false},
		{"authorized user with empty name", true, "", false},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			store := sess.InitStoreFromTimeNow()

			r := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
			w := httptest.NewRecorder()

			session, _ := sess.GetSessionFromStore(store, r)
			session.Values[config.SessionAuthorized] = tt.authorized
			session.Values[config.SessionUserName] = tt.username
			session.Save(r, w)

			h := &executedHandler{}

			AuthOnly(store)(h).ServeHTTP(w, r)

			assert.Equal(t, h.IsExecuted, tt.shouldExecute)
		})
	}
}
