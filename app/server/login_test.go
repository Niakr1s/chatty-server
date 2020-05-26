package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/niakr1s/chatty-server/app/constants"
	"github.com/niakr1s/chatty-server/app/internal/sess"
	"github.com/niakr1s/chatty-server/app/internal/validator"
	"github.com/niakr1s/chatty-server/app/models"
	"github.com/stretchr/testify/assert"
)

func TestServer_Login(t *testing.T) {
	s := NewMemoryServer()

	username := "user"

	u := models.User{Name: username}
	b, _ := json.Marshal(u)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))

	s.Login(w, r)
	assert.Equal(t, http.StatusOK, w.Code)

	loggedU, err := s.dbStore.LoggedDB.Get(username)
	assert.NoError(t, err)

	err = validator.Validate.Struct(loggedU)
	assert.NoError(t, err)
	assert.Equal(t, username, loggedU.Name)

	session, _ := sess.GetSessionFromStore(s.cookieStore, r)

	assert.Equal(t, loggedU.Name, session.Values[constants.SessionUserName].(string))
	assert.Equal(t, loggedU.LoginToken, session.Values[constants.SessionLoginToken].(string))
}

func TestServer_LoginSameUserTwice(t *testing.T) {
	s := NewMemoryServer()

	username := "user"

	u := models.User{Name: username}
	b, _ := json.Marshal(u)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))
	s.Login(w, r)
	assert.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))
	s.Login(w, r)
	assert.NotEqual(t, http.StatusOK, w.Code)
}
