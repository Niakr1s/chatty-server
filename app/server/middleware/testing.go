package middleware

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mockHandlerFail(t *testing.T) http.Handler {
	t.Helper()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Fail(t, "this code shouldn't exec")
	})
}

type executedHandler struct {
	IsExecuted bool
}

func (h *executedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.IsExecuted = true
}
