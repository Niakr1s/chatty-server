package events

import (
	"server2/app/er"
	"server2/app/models"
)

// MessageEvent ...
type MessageEvent struct {
	*models.Message
}

// NewMessageEvent ...
func NewMessageEvent(msg *models.Message) *MessageEvent {
	return &MessageEvent{msg}
}

// InChat ...
func (e *MessageEvent) InChat() (string, error) {
	if e.Chat == "" {
		return "", er.ErrGlobalEvent
	}
	return e.Chat, nil
}
