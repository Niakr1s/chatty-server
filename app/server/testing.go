package server

import (
	"context"
	"net/http"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/niakr1s/chatty-server/app/config"
	"github.com/niakr1s/chatty-server/app/models"
)

func mockUser(t *testing.T) models.User {
	t.Helper()
	return models.User{Name: "user", Password: "password"}
}

func appendSessionToRequest(t *testing.T, r *http.Request, s *sessions.Session) *http.Request {
	t.Helper()
	r = r.WithContext(context.WithValue(r.Context(), config.CtxSessionKey, s))
	return r
}
