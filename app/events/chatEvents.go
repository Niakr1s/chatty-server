package events

import (
	"fmt"
	"time"

	"github.com/niakr1s/chatty-server/app/models"
)

// ChatEvent ...
type ChatEvent struct {
	models.Chat
	Time models.UnixTime
}

// NewChatEvent ...
func NewChatEvent(chatname string, t time.Time) *ChatEvent {
	return &ChatEvent{Chat: models.Chat{ChatName: chatname}, Time: models.UnixTime(t)}
}

func (ce *ChatEvent) String() string {
	return fmt.Sprintf("chat: %s, time: %v", ce.ChatName, ce.Time)
}

// ChatCreatedEvent ...
type ChatCreatedEvent struct {
	*ChatEvent
}

// NewChatCreatedEvent ...
func NewChatCreatedEvent(chatname string, t time.Time) *ChatCreatedEvent {
	return &ChatCreatedEvent{ChatEvent: NewChatEvent(chatname, t)}
}

func (ce *ChatCreatedEvent) String() string {
	return fmt.Sprintf("chat created: %v", ce.ChatEvent)
}

// ChatRemovedEvent ...
type ChatRemovedEvent struct {
	*ChatEvent
}

// NewChatRemovedEvent ...
func NewChatRemovedEvent(chatname string, t time.Time) *ChatRemovedEvent {
	return &ChatRemovedEvent{ChatEvent: NewChatEvent(chatname, t)}
}

func (ce *ChatRemovedEvent) String() string {
	return fmt.Sprintf("chat removed: %v", ce.ChatEvent)
}
