package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/niakr1s/chatty-server/app/config"
	"github.com/niakr1s/chatty-server/app/server/sess"
	"github.com/stretchr/testify/assert"
)

func TestAddSessionToContext(t *testing.T) {
	store := sess.InitStoreFromTimeNow()

	h := func(w http.ResponseWriter, r *http.Request) {
		s := r.Context().Value(config.CtxSessionKey).(*sessions.Session)
		assert.Equal(t, s.Values["test"], "test")
	}

	r := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(""))
	w := httptest.NewRecorder()

	s, _ := sess.GetSessionFromStore(store, r)
	s.Values["test"] = "test"
	s.Save(r, w)

	AddSessionToContext(store)(http.HandlerFunc(h)).ServeHTTP(w, r)
}
