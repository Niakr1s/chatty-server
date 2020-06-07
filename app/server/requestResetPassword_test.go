package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/niakr1s/chatty-server/app/email"
	"github.com/niakr1s/chatty-server/app/models"
	"github.com/stretchr/testify/assert"
)

func TestServer_RequestResetPassword(t *testing.T) {
	tests := []struct {
		name       string
		username   string
		okExpected bool
	}{
		{"valid", mockUsername, true},
		{"not registered", "not registered user", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newMockServer()
			m := email.NewMockMailer()
			s.mailer = m

			s.dbStore.UserDB.Store(models.NewFullUser(mockUsername, mockEmail, mockPassword))

			u := models.NewUser(tt.username)
			b, _ := json.Marshal(u)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))

			s.RequestResetPassword(w, r)
			assert.Equal(t, tt.okExpected, w.Code == http.StatusOK)

			if !tt.okExpected {
				return
			}

			fullU, _ := s.dbStore.UserDB.Get(tt.username)
			assert.Equal(t, fullU.PasswordResetToken, m.ResetPasswordToken)
		})
	}
}
