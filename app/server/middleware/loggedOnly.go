package middleware

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/niakr1s/chatty-server/app/db"
	"github.com/niakr1s/chatty-server/app/er"
	"github.com/niakr1s/chatty-server/app/internal/httputil"
	"github.com/niakr1s/chatty-server/app/internal/sess"
)

// LoggedOnly rejects user with invalid loginToken
// Stores username in context
func LoggedOnly(s sessions.Store, loggedDB db.LoggedDB) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := sess.GetSessionFromStore(s, r)
			if err != nil {
				httputil.WriteSessionError(w)
				return
			}

			if !sess.IsLogged(session, loggedDB) {
				httputil.WriteError(w, er.ErrNotLogged, http.StatusUnauthorized)
				return
			}

			username, err := sess.GetUserName(session)
			if err != nil {
				httputil.WriteSessionError(w)
				return
			}

			r = r.WithContext(sess.SetUserNameIntoCtx(r.Context(), username))

			h.ServeHTTP(w, r)
		})
	}
}
