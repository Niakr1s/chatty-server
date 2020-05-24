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
	t.Run("unathorized", func(t *testing.T) {
		h := mockHandlerFail(t)
		r := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
		w := httptest.NewRecorder()
		AuthOnly(h).ServeHTTP(w, r)
	})

	t.Run("authorized", func(t *testing.T) {
		store := sess.InitStoreFromTimeNow()

		r := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
		w := httptest.NewRecorder()

		session, _ := sess.GetSessionFromStore(store, r)
		session.Values[config.SessionAuthorized] = true
		session.Save(r, w)

		r = r.WithContext(sess.ContextWithSession(r.Context(), session))

		h := &executedHandler{}

		AuthOnly(h).ServeHTTP(w, r)

		assert.True(t, h.IsExecuted)
	})
}
