package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/niakr1s/chatty-server/app/internal/sess"
	"github.com/stretchr/testify/assert"
)

func TestServer_AuthLogin(t *testing.T) {
	s := newMockServer()

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
		// {"empty username", "", false}, // emtpty user should be handled with LoggedOnly middleware
	}
	for _, tt := range tests {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))

		r = r.WithContext(sess.SetUserNameIntoCtx(r.Context(), tt.username))

		s.AuthLogin(w, r)

		assert.Equal(t, tt.okExpected, w.Code == http.StatusOK)
	}
}
