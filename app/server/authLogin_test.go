package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/niakr1s/chatty-server/app/config"
	"github.com/stretchr/testify/assert"
)

func TestServer_AuthLogin(t *testing.T) {
	s := NewMemoryServer()

	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name       string
		username   string
		okExpected bool
	}{
		{"valid username", "username", true},
		{"empty username", "", false},
	}
	for _, tt := range tests {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
		if tt.username != "" {
			r = r.WithContext(context.WithValue(r.Context(), config.CtxUserNameKey, tt.username))
		}

		s.AuthLogin(w, r)

		assert.Equal(t, tt.okExpected, w.Code == http.StatusOK)
	}
}
