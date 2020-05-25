package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/niakr1s/chatty-server/app/config"
	"github.com/niakr1s/chatty-server/app/db/logged"
	"github.com/niakr1s/chatty-server/app/models"
	"github.com/niakr1s/chatty-server/app/server/sess"
	"github.com/stretchr/testify/assert"
)

func TestLoggedOnly_NewUser(t *testing.T) {
	store := sess.InitStoreFromTimeNow()
	loggedDB := logged.NewMemoryDB()
	h := &executedHandler{}

	username := "user"

	u := models.User{Name: username}
	b, _ := json.Marshal(u)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))

	LoggedOnly(store, loggedDB)(h).ServeHTTP(w, r)

	assert.NotEqual(t, http.StatusOK, w.Code)
	assert.False(t, h.IsExecuted)
}

func TestLoggedOnly_LoggedUser(t *testing.T) {
	username := "user"

	store := sess.InitStoreFromTimeNow()
	loggedDB := logged.NewMemoryDB()
	h := &executedHandler{}

	u, _ := loggedDB.Login(username)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))

	session, _ := sess.GetSessionFromStore(store, r)
	session.Values[config.SessionUserName] = u.Name
	session.Values[config.SessionLoginToken] = u.LoginToken
	session.Save(r, w)

	LoggedOnly(store, loggedDB)(h).ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, h.IsExecuted)
}
