package pool

import (
	"server2/app/pool/events"
	"server2/app/store"
)

// FilterPassIfUserInChat ...
func FilterPassIfUserInChat(store *store.Store, username string) events.FilterPass {
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
