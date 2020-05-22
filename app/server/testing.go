package server

import (
	"server2/app/models"
	"testing"
)

func mockUser(t *testing.T) models.User {
	t.Helper()
	return models.User{Name: "user", Password: "password"}
}
