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
		b, _ := json.Marshal(&tt.user)

		r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(b))
		w := httptest.NewRecorder()

		s.Register(w, r)

		assert.Equal(t, (w.Code == http.StatusOK), tt.okExpected)
	}
}
