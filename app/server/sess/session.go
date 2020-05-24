package sess

import (
	"context"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/niakr1s/chatty-server/app/config"
	"github.com/niakr1s/chatty-server/app/er"
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

// IsAuthorized ...
func IsAuthorized(session *sessions.Session) bool {
	v := session.Values[config.SessionAuthorized]
	if v == nil {
		return false
	}

	res, ok := v.(bool)
	if !ok {
		return false
	}

	return res
}

// GetUserName ...
func GetUserName(session *sessions.Session) (string, error) {
	v := session.Values[config.SessionUserName]
	if v == nil {
		return "", er.ErrUserNameIsEmpty
	}

	res, ok := v.(string)
	if !ok {
		return "", er.ErrConvertType
	}

	if res == "" {
		return "", er.ErrUserNameIsEmpty
	}

	return res, nil
}
