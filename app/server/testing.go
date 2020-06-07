package server

import (
	"github.com/niakr1s/chatty-server/app/db"
	"github.com/niakr1s/chatty-server/app/db/memory"
	"github.com/niakr1s/chatty-server/app/email"
)

const (
	mockUsername = "user"
	mockEmail    = "email@example.org"
	mockPassword = "password"
)

// newMockServer ...
func newMockServer() *Server {
	u := memory.NewUserDB()
	c := memory.NewChatDB()
	l := memory.NewLoggedDB()
	m := memory.NewMessageDB()
	return NewServer(db.NewStore(u, c, l, m), email.NewMockMailer())
}
