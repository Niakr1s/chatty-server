package models

import (
	"fmt"
	"time"

	"github.com/niakr1s/chatty-server/app/er"
	"github.com/niakr1s/chatty-server/app/internal/validator"
)

// Message is representation of message.
// Time must contain UTC time, not local.
type Message struct {
	User
	Chat
	ID   int      `json:"id"`
	Text string   `json:"text" validate:"required"`
	Time UnixTime `json:"time"`
}

// NewMessage constructs message with Time=time.Now().UTC()
func NewMessage(username, text, chat string) *Message {
	return &Message{User: User{UserName: username},
		Text: text, Chat: Chat{ChatName: chat},
		Time: UnixTime(time.Now().UTC())}
}

// WithTime ...
func (m *Message) WithTime(t time.Time) *Message {
	m.Time = UnixTime(t)
	return m
}

// ValidateBeforeStoring ...
func (m *Message) ValidateBeforeStoring() error {
	if time.Now().UTC().Sub(time.Time(m.Time)) > time.Minute {
		return er.ErrTooOld
	}
	return validator.Validate.Struct(*m)
}

func (m *Message) String() string {
	return fmt.Sprintf("chat: %s, id: %d, user: %s, text: %s, time: %v", m.Chat, m.ID, m.UserName, m.Text, m.Time)
}
