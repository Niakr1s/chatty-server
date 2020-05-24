package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/niakr1s/chatty-server/app/config"
	"github.com/niakr1s/chatty-server/app/models"
	"github.com/niakr1s/chatty-server/app/server/sess"

	"github.com/stretchr/testify/assert"
)

func TestServer_Register(t *testing.T) {
	s := NewMemoryServer()

	tests := []struct {
		name       string
		user       models.User
		okExpected bool
	}{
		{
			"valid user",
			mockUser(t),
			true,
		},
		{
			"user with empty password",
			models.User{Name: "user"},
			false,
		},
		{
			"user with empty name",
			models.User{Password: "password"},
			false,
		},
	}
	for _, tt := range tests {
		store := sess.InitStoreFromTimeNow()

		b, _ := json.Marshal(&tt.user)

		r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(b))

		session, _ := sess.GetSessionFromStore(store, r)
		r = appendSessionToRequest(t, r, session)

		w := httptest.NewRecorder()

		s.Register(w, r)

		assert.Equal(t, (w.Code == http.StatusOK), tt.okExpected)

		if authorized := sess.GetSessionFromContext(r.Context()).Values[config.SessionAuthorized]; tt.okExpected {
			assert.Equal(t, authorized, tt.okExpected)
		} else {
			assert.Nil(t, authorized)
		}
	}
}
