package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/niakr1s/chatty-server/app/constants"
	"github.com/niakr1s/chatty-server/app/db/memory"
	"github.com/niakr1s/chatty-server/app/internal/sess"
	"github.com/niakr1s/chatty-server/app/models"
	"github.com/stretchr/testify/assert"
)

func TestAdminOnly(t *testing.T) {
	tests := []struct {
		name         string
		status       models.UserStatus
		execExpected bool
	}{
		{"new user", models.UserStatus{}, false},
		{"unverified admin", models.UserStatus{Admin: true, Verified: false}, false},
		{"verified no-admin", models.UserStatus{Admin: false, Verified: true}, false},
		{"verified admin", models.UserStatus{Admin: true, Verified: true}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			username := "user"

			store := sess.InitStoreFromTimeNow()
			loggedDB := memory.NewLoggedDB()
			h := &executedHandler{}

			u, _ := loggedDB.Login(username)
			u.UserStatus = tt.status
			loggedDB.Update(u)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))

			session, _ := sess.GetSessionFromStore(store, r)
			session.Values[constants.SessionUserName] = u.UserName
			session.Values[constants.SessionLoginToken] = u.LoginToken
			session.Save(r, w)

			AdminOnly(store, loggedDB)(h).ServeHTTP(w, r)

			assert.Equal(t, tt.execExpected, h.IsExecuted)
		})
	}
}
