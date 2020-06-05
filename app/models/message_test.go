package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	username = "user"
	text     = "text"
	chat     = "chat"
)

func TestMessage_ValidateBeforeStoring(t *testing.T) {
	tests := []struct {
		name    string
		m       *Message
		wantErr bool
	}{
		{
			"valid message",
			NewMessage(username, text, chat).WithTime(time.Now().UTC()),
			false,
		},
		{
			"old message",
			NewMessage(username, text, chat).WithTime(time.Now().UTC().Add(-time.Hour * 24)),
			true,
		},
		{
			"empty user",
			NewMessage("", text, chat).WithTime(time.Now().UTC()),
			true,
		},
		{
			"empty text",
			NewMessage(username, "", chat).WithTime(time.Now().UTC()),
			true,
		},
		{
			"empty chat",
			NewMessage(username, text, "").WithTime(time.Now().UTC()),
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.m.ValidateBeforeStoring()
			assert.Equal(t, err != nil, tt.wantErr)
		})
	}
}
