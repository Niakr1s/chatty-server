package events

import (
	"fmt"

	"github.com/niakr1s/chatty-server/app/models"
)

// SystemMessageEvent ...
type SystemMessageEvent struct {
	models.Chat
	models.User
}

func newSystemMessageEvent(chatname, username string) *SystemMessageEvent {
	return &SystemMessageEvent{Chat: models.NewChat(chatname), User: models.NewUser(username)}
}

func (e *SystemMessageEvent) String() string {
	return fmt.Sprintf("chat: %s, user: %s", e.ChatName, e.UserName)
}

// SystemMessageChatJoinEvent ...
type SystemMessageChatJoinEvent struct {
	*SystemMessageEvent
}

// NewSystemMessageChatJoinEvent ...
func NewSystemMessageChatJoinEvent(chatname, username string) *SystemMessageChatJoinEvent {
	return &SystemMessageChatJoinEvent{SystemMessageEvent: newSystemMessageEvent(chatname, username)}
}

func (e *SystemMessageChatJoinEvent) String() string {
	return fmt.Sprintf("sysMsg: chat join: %v", e.SystemMessageEvent)
}

// SystemMessageChatLeaveEvent ...
type SystemMessageChatLeaveEvent struct {
	*SystemMessageEvent
}

// NewSystemMessageChatLeaveEvent ...
func NewSystemMessageChatLeaveEvent(chatname, username string) *SystemMessageChatLeaveEvent {
	return &SystemMessageChatLeaveEvent{SystemMessageEvent: newSystemMessageEvent(chatname, username)}
}

func (e *SystemMessageChatLeaveEvent) String() string {
	return fmt.Sprintf("sysMsg: chat leave: %v", e.SystemMessageEvent)
}
