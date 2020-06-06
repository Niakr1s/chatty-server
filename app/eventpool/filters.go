package eventpool

import (
	"github.com/niakr1s/chatty-server/app/db"
	"github.com/niakr1s/chatty-server/app/er"
	"github.com/niakr1s/chatty-server/app/events"
)

// FilterPass should return true if event is passable
type FilterPass func(e events.Event) bool

// FilterPassAlways event is always passable
func FilterPassAlways(e events.Event) bool {
	return true
}

// FilterPassLogoutEvents passes only LogoutEvents
func FilterPassLogoutEvents(e events.Event) bool {
	switch e.(type) {
	case *events.LogoutEvent:
		return true
	default:
		return false
	}
}

// FilterPassIfUserInChat event is passable only if user is in same chat as event occurs in
// or if event global
func FilterPassIfUserInChat(chatDB db.ChatDB, username string) FilterPass {
	return FilterPass(func(e events.Event) bool {
		inChat, err := e.InChat()
		if err != nil {
			return err == er.ErrGlobalEvent
		}

		return chatDB.IsInChat(inChat, username)
	})
}
