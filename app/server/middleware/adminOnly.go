package middleware

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/niakr1s/chatty-server/app/db"
	"github.com/niakr1s/chatty-server/app/er"
	"github.com/niakr1s/chatty-server/app/internal/httputil"
	"github.com/niakr1s/chatty-server/app/internal/sess"
)

// AdminOnly rejects user with invalid loginToken
// Stores username in context
func AdminOnly(s sessions.Store, loggedDB db.LoggedDB) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, err := sess.GetSessionFromStore(s, r)
			if err != nil {
				httputil.WriteSessionError(w)
				return
			}

			username, err := sess.GetUserName(session)
			if err != nil {
				httputil.WriteSessionError(w)
				return
			}

			loggedDB.Lock()
			loggedU, err := loggedDB.Get(username)
			if err != nil {
				httputil.WriteError(w, er.ErrNotLogged, http.StatusBadRequest)
				return
			}
			loggedDB.Unlock()

			if loggedU.Verified && loggedU.Admin {
				r = r.WithContext(sess.SetUserNameIntoCtx(r.Context(), username))
				h.ServeHTTP(w, r)
			}
		})
	}
}
