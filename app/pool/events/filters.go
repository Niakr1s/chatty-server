package events

import "server2/app/store"

// FilterPass ...
type FilterPass func(e Event) bool

// FilterPassIfUserInChat ...
func FilterPassIfUserInChat(store *store.Store, username string) FilterPass {
	return FilterPass(func(e Event) bool {
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

// FilterPassAlways ...
func FilterPassAlways(e Event) bool {
	return true
}
