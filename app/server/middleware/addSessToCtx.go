package middleware

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/niakr1s/chatty-server/app/server/sess"
	"github.com/niakr1s/chatty-server/app/server/util"
)

// AddSessionToContext ...
func AddSessionToContext(store *sessions.CookieStore) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := sess.GetSessionFromStore(store, r)
			if err != nil {
				util.WriteError(w, err, http.StatusInternalServerError)
				return
			}

			r = r.WithContext(sess.ContextWithSession(r.Context(), session))

			h.ServeHTTP(w, r)
		})
	}
}
