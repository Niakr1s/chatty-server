package events

import (
	"server2/app/er"
	"time"
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

// ChatCreatedEvent ...
type ChatCreatedEvent struct {
	*ChatEvent
}

// NewChatCreatedEvent ...
func NewChatCreatedEvent(chatname string, t time.Time) *ChatCreatedEvent {
	return &ChatCreatedEvent{ChatEvent: NewChatEvent(chatname, t)}
}

// ChatRemovedEvent ...
type ChatRemovedEvent struct {
	*ChatEvent
}

// NewChatRemovedEvent ...
func NewChatRemovedEvent(chatname string, t time.Time) *ChatRemovedEvent {
	return &ChatRemovedEvent{ChatEvent: NewChatEvent(chatname, t)}
}
