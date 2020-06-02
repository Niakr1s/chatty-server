package server

import (
	"github.com/niakr1s/chatty-server/app/db"
	"github.com/niakr1s/chatty-server/app/db/chat"
	"github.com/niakr1s/chatty-server/app/db/logged"
	"github.com/niakr1s/chatty-server/app/db/message"
	"github.com/niakr1s/chatty-server/app/db/user"
	"github.com/niakr1s/chatty-server/app/email"
)

// newMockServer ...
func newMockServer() *Server {
	u := user.NewMemoryDB()
	c := chat.NewMemoryDB()
	l := logged.NewMemoryDB()
	m := message.NewMemoryDB()
	return NewServer(db.NewStore(u, c, l, m), email.NewMockMailer())
}
