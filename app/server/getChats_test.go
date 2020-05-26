package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer_GetChats(t *testing.T) {
	tests := []struct {
		name       string
		okExpected bool
	}{
		{
			"valid",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewMemoryServer()
			s.dbStore.ChatDB.Add("chat")

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))

			s.GetChats(w, r)

			assert.Equal(t, tt.okExpected, w.Code == http.StatusOK)
		})
	}
}
