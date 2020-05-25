package events

import (
	"fmt"
	"time"

	"github.com/niakr1s/chatty-server/app/er"
)

// ChatEvent ...
type ChatEvent struct {
	Chatname string
	Time     time.Time
}

// NewChatEvent ...
func NewChatEvent(chatname string, t time.Time) *ChatEvent {
	return &ChatEvent{Chatname: chatname, Time: t}
}

// InChat ...
func (ce *ChatEvent) InChat() (string, error) {
	if ce.Chatname == "" {
		return "", er.ErrGlobalEvent
	}
	return ce.Chatname, nil
}

func (ce *ChatEvent) String() string {
	return fmt.Sprintf("chat: %s, time: %v", ce.Chatname, ce.Time)
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
