package events

import (
	"github.com/niakr1s/chatty-server/app/er"
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

// InChat ...
func (e *MessageEvent) InChat() (string, error) {
	if e.Chat == "" {
		return "", er.ErrGlobalEvent
	}
	return e.Chat, nil
}
