package server

import (
	"testing"

	"github.com/niakr1s/chatty-server/app/models"
)

func mockUser(t *testing.T) models.User {
	t.Helper()
	return models.User{Name: "user", Password: "password"}
}
