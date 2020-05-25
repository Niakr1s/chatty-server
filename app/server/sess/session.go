package sess

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/niakr1s/chatty-server/app/config"
	"github.com/niakr1s/chatty-server/app/db"
	"github.com/niakr1s/chatty-server/app/er"
)

// GetSessionFromStore ...
func GetSessionFromStore(store *sessions.CookieStore, r *http.Request) (*sessions.Session, error) {
	return store.Get(r, config.SessionName)
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

// IsLogged ...
func IsLogged(session *sessions.Session, loggedDB db.LoggedDB) bool {
	usernameI := session.Values[config.SessionUserName]
	loginTokenI := session.Values[config.SessionLoginToken]
	if usernameI == nil || loginTokenI == nil {
		return false
	}

	username, usernameOk := usernameI.(string)
	loginToken, loginTokenOk := loginTokenI.(string)

	if !usernameOk || !loginTokenOk {
		return false
	}

	loggedDB.Lock()
	defer loggedDB.Unlock()
	loggedU, err := loggedDB.Get(username)
	if err != nil {
		return false
	}

	return loggedU.LoginToken == loginToken
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
