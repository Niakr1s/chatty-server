package sess

import (
	"context"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/niakr1s/chatty-server/app/config"
)

// GetSessionFromStore ...
func GetSessionFromStore(store *sessions.CookieStore, r *http.Request) (*sessions.Session, error) {
	return store.Get(r, config.SessionName)
}

// GetSessionFromContext ...
func GetSessionFromContext(ctx context.Context) *sessions.Session {
	return ctx.Value(config.CtxSessionKey).(*sessions.Session)
}
