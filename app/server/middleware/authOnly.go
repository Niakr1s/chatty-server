package middleware

import (
	"net/http"

	"github.com/niakr1s/chatty-server/app/config"
	"github.com/niakr1s/chatty-server/app/er"
	"github.com/niakr1s/chatty-server/app/server/sess"
	"github.com/niakr1s/chatty-server/app/server/util"
)

// AuthOnly ...
func AuthOnly(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := sess.GetSessionFromContext(r.Context())

		if session == nil {
			util.WriteError(w, er.ErrUnathorized, http.StatusUnauthorized)
			return
		}

		auth := session.Values[config.SessionAuthorized]

		if auth == nil || auth == false {
			return
		}

		h.ServeHTTP(w, r)
	})
}
