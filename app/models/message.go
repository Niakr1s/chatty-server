package models

import (
	"fmt"
	"time"

	"github.com/niakr1s/chatty-server/app/er"
	"github.com/niakr1s/chatty-server/app/internal/validator"
)

// Message ...
type Message struct {
	ID       int      `json:"id"`
	Username string   `json:"username" validate:"required"`
	Text     string   `json:"text" validate:"required"`
	Chat     string   `json:"chat" validate:"required"`
	Time     UnixTime `json:"time"`
}

// NewMessage ...
func NewMessage(username, text, chat string) *Message {
	return &Message{Username: username, Text: text, Chat: chat, Time: UnixTime(time.Now())}
}

// WithTime ...
func (m *Message) WithTime(t time.Time) *Message {
	m.Time = UnixTime(t)
	return m
}

// ValidateBeforeStoring ...
func (m *Message) ValidateBeforeStoring() error {
	if time.Now().Sub(time.Time(m.Time)) > time.Minute {
		return er.ErrTooOld
	}
	return validator.Validate.Struct(*m)
}

func (m *Message) String() string {
	return fmt.Sprintf("chat: %s, id: %d, user: %s, text: %s, time: %v", m.Chat, m.ID, m.Username, m.Text, m.Time)
}
