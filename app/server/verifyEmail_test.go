package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/niakr1s/chatty-server/app/email"
	"github.com/niakr1s/chatty-server/app/models"
	"github.com/stretchr/testify/assert"
)

func TestServer_VerifyEmail(t *testing.T) {
	s := newMockServer()
	m := s.mailer.(*email.MockMailer)

	u := models.NewFullUser("user", "user@example.com", "password")
	b, _ := json.Marshal(&u)
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))
	w := httptest.NewRecorder()

	s.Register(w, r)

	r = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
	w = httptest.NewRecorder()
	r = mux.SetURLVars(r, map[string]string{"username": m.User, "activationToken": m.ActivationToken})

	s.VerifyEmail(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}
