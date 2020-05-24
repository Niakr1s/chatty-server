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
	session := ctx.Value(config.CtxSessionKey)

	if session == nil {
		return nil
	}

	return session.(*sessions.Session)
}

// ContextWithSession ...
func ContextWithSession(ctx context.Context, session *sessions.Session) context.Context {
	return context.WithValue(ctx, config.CtxSessionKey, session)
}
