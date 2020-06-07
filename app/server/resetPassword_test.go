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

func TestServer_ResetPassword(t *testing.T) {
	const newPassword = "newpassword"
	const mockPasswordHash = "12345"

	tests := []struct {
		name       string
		username   string
		resetToken string
		okExpected bool
	}{
		{"valid", mockUsername, mockToken, true},
		{"valid vith invalid token", mockUsername, "invalid token", false},
		{"valid vith empty token", mockUsername, "", false},
		{"empty user with token from valid user", "", mockToken, false},
		{"empty user with empty token", "", "", false},
		{"non-registered user with token from valid user", "non-registered-user", mockToken, false},
		{"non-registered user", "non-registered-user", "2395870", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newMockServer()
			m := email.NewMockMailer()
			s.mailer = m

			storedU := models.NewFullUser(mockUsername, mockEmail, mockPassword)
			storedU.PasswordResetToken = mockToken
			storedU.PasswordHash = mockPasswordHash
			storedU.ActivationToken = mockToken
			s.dbStore.UserDB.Store(storedU)

			type userWithPass struct {
				models.User
				models.Pass
			}
			u := userWithPass{User: models.NewUser(tt.username), Pass: models.Pass{Password: newPassword, PasswordResetToken: tt.resetToken}}
			b, _ := json.Marshal(u)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))

			s.ResetPassword(w, r)
			assert.Equal(t, tt.okExpected, w.Code == http.StatusOK)
			if !tt.okExpected {
				return
			}

			gotU, _ := s.dbStore.UserDB.Get(tt.username)
			assert.Equal(t, mockUsername, gotU.UserName)
			assert.Empty(t, gotU.PasswordResetToken)
			assert.NotEmpty(t, gotU.PasswordHash)
			assert.NotEqual(t, storedU.PasswordHash, gotU.PasswordHash)
		})
	}
}
