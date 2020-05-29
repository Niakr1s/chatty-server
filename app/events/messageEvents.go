package events

import (
	"github.com/niakr1s/chatty-server/app/models"
)

// MessageEvent ...
type MessageEvent struct {
	*models.Message
}

// NewMessageEvent ...
func NewMessageEvent(msg *models.Message) *MessageEvent {
	return &MessageEvent{msg}
}

func (e *MessageEvent) String() string {
	return e.Message.String()
}
