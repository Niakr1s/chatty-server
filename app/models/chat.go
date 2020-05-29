package models

import "github.com/niakr1s/chatty-server/app/er"

// Chat ...
type Chat struct {
	ChatName string `json:"chat" validate:"required"`
}

// NewChat ...
func NewChat(chatname string) Chat {
	return Chat{ChatName: chatname}
}

// InChat ...
func (c *Chat) InChat() (string, error) {
	if c.ChatName == "" {
		return "", er.ErrGlobalEvent
	}
	return c.ChatName, nil
}
