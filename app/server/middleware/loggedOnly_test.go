package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/niakr1s/chatty-server/app/constants"
	"github.com/niakr1s/chatty-server/app/db/memory"
	"github.com/niakr1s/chatty-server/app/internal/sess"
	"github.com/niakr1s/chatty-server/app/models"
	"github.com/stretchr/testify/assert"
)

func TestLoggedOnly_NewUser(t *testing.T) {
	store := sess.InitStoreFromTimeNow()
	loggedDB := memory.NewLoggedDB()
	h := &executedHandler{}

	username := "user"

	u := models.User{UserName: username}
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
	loggedDB := memory.NewLoggedDB()
	h := &executedHandler{}

	u, _ := loggedDB.Login(username)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))

	session, _ := sess.GetSessionFromStore(store, r)
	session.Values[constants.SessionUserName] = u.UserName
	session.Values[constants.SessionLoginToken] = u.LoginToken
	session.Save(r, w)

	LoggedOnly(store, loggedDB)(h).ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, h.IsExecuted)
}
