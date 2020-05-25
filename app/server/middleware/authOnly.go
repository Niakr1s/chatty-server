package middleware

import (
	"context"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/niakr1s/chatty-server/app/config"
	"github.com/niakr1s/chatty-server/app/er"
	"github.com/niakr1s/chatty-server/app/server/httputil"
	"github.com/niakr1s/chatty-server/app/server/sess"
)

// AuthOnly reject unauthorized user
// stores username in context
func AuthOnly(s sessions.Store) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := sess.GetSessionFromStore(s, r)
			if err != nil {
				httputil.WriteSessionError(w)
				return
			}

			if session == nil {
				httputil.WriteError(w, er.ErrUnathorized, http.StatusUnauthorized)
				return
			}

			if !sess.IsAuthorized(session) {
				httputil.WriteError(w, er.ErrUnathorized, http.StatusUnauthorized)
				return
			}

			username, err := sess.GetUserName(session)
			if err != nil || username == "" {
				httputil.WriteError(w, er.ErrUserNameIsEmpty, http.StatusUnauthorized)
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), config.CtxUserNameKey, username))

			h.ServeHTTP(w, r)
		})
	}
}
