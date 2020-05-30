package events

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewEventWithType(t *testing.T) {
	const (
		user = "user"
		chat = "chat"
	)
	testCases := []struct {
		e            Event
		expectedType string
	}{
		{NewLoginEvent(user, chat, time.Now()), "LoginEvent"},
		{NewLogoutEvent(user, chat, time.Now()), "LogoutEvent"},
		{NewChatEvent(chat, time.Now()), "ChatEvent"},
	}

	for _, tt := range testCases {
		ewt := NewEventWithType(tt.e)
		assert.Equal(t, tt.expectedType, ewt.Type)
	}
}
