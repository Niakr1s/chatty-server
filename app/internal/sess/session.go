package sess

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/niakr1s/chatty-server/app/constants"
	"github.com/niakr1s/chatty-server/app/db"
	"github.com/niakr1s/chatty-server/app/er"
)

// GetSessionFromStore ...
func GetSessionFromStore(store sessions.Store, r *http.Request) (*sessions.Session, error) {
	return getSessionFromStore(store, r, 0)
}

// GetSessionFromStore deletes cookies and tries again one time
func getSessionFromStore(store sessions.Store, r *http.Request, times int) (*sessions.Session, error) {
	if times > 1 {
		return nil, er.ErrSession
	}
	session, err := store.Get(r, constants.SessionName)
	if err != nil {
		session, err = store.New(r, constants.SessionName)
		if err != nil {
			r.Header.Set("Cookie", "")
			return getSessionFromStore(store, r, times+1)
		}
	}
	return session, nil
}

// IsLogged ...
func IsLogged(session *sessions.Session, loggedDB db.LoggedDB) bool {
	usernameI := session.Values[constants.SessionUserName]
	loginTokenI := session.Values[constants.SessionLoginToken]
	if usernameI == nil || loginTokenI == nil {
		return false
	}

	username, usernameOk := usernameI.(string)
	loginToken, loginTokenOk := loginTokenI.(string)

	if !usernameOk || !loginTokenOk {
		return false
	}

	loggedU, err := loggedDB.Get(username)
	if err != nil {
		return false
	}

	return loggedU.LoginToken == loginToken
}

// GetUserName ...
func GetUserName(session *sessions.Session) (string, error) {
	v := session.Values[constants.SessionUserName]
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
