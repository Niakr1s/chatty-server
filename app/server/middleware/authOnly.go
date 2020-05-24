package middleware

import (
	"context"
	"net/http"

	"github.com/niakr1s/chatty-server/app/config"
	"github.com/niakr1s/chatty-server/app/er"
	"github.com/niakr1s/chatty-server/app/server/httputil"
	"github.com/niakr1s/chatty-server/app/server/sess"
)

// AuthOnly ...
func AuthOnly(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := sess.GetSessionFromContext(r.Context())

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
