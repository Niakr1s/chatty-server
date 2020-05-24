package pool

import (
	"github.com/niakr1s/chatty-server/app/db"
	"github.com/niakr1s/chatty-server/app/pool/events"
)

// FilterPassIfUserInChat ...
func FilterPassIfUserInChat(store *db.Store, username string) events.FilterPass {
	return events.FilterPass(func(e events.Event) bool {
		inChat, err := e.InChat()

		if err != nil {
			return false
		}

		chat, err := store.ChatDB.Get(inChat)

		if err != nil {
			return false
		}

		return chat.IsInChat(username)
	})
}
