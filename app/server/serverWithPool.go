package server

import (
	"github.com/niakr1s/chatty-server/app/db/chat"
	"github.com/niakr1s/chatty-server/app/db/logged"
	"github.com/niakr1s/chatty-server/app/pool"
	"github.com/niakr1s/chatty-server/app/pool/events"
)

// WithPool ...
type WithPool struct {
	*Server
	pool *pool.Pool
}

// NewServerWithPool ...
func NewServerWithPool(s *Server) *WithPool {
	res := &WithPool{Server: s, pool: pool.NewPool()}
	ch := res.pool.GetInputChan()
	res.Store.LoggedDB = logged.NewNotifyDB(res.Store.LoggedDB, ch)
	res.Store.ChatDB = chat.NewNotifyDB(res.Store.ChatDB, ch)
	res.pool = res.pool.WithUserChFilter(func(username string) events.FilterPass {
		return pool.FilterPassIfUserInChat(res.Store, username)
	})
	res.pool.Run()
	return res
}
