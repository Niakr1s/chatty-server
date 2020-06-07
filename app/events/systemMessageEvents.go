package events

import (
	"fmt"

	"github.com/niakr1s/chatty-server/app/models"
)

// SystemMessageEvent ...
type SystemMessageEvent struct {
	models.Chat
	models.User
	Time models.UnixTime `json:"time"`
}

func newSystemMessageEvent(chatname, username string, t models.UnixTime) *SystemMessageEvent {
	return &SystemMessageEvent{Chat: models.NewChat(chatname), User: models.NewUser(username), Time: t}
}

func (e *SystemMessageEvent) String() string {
	return fmt.Sprintf("chat: %s, user: %s", e.ChatName, e.UserName)
}

// SystemMessageChatJoinEvent ...
type SystemMessageChatJoinEvent struct {
	*SystemMessageEvent
}

// NewSystemMessageChatJoinEvent ...
func NewSystemMessageChatJoinEvent(chatname, username string, t models.UnixTime) *SystemMessageChatJoinEvent {
	return &SystemMessageChatJoinEvent{SystemMessageEvent: newSystemMessageEvent(chatname, username, t)}
}

func (e *SystemMessageChatJoinEvent) String() string {
	return fmt.Sprintf("sysMsg: chat join: %v", e.SystemMessageEvent)
}

// SystemMessageChatLeaveEvent ...
type SystemMessageChatLeaveEvent struct {
	*SystemMessageEvent
}

// NewSystemMessageChatLeaveEvent ...
func NewSystemMessageChatLeaveEvent(chatname, username string, t models.UnixTime) *SystemMessageChatLeaveEvent {
	return &SystemMessageChatLeaveEvent{SystemMessageEvent: newSystemMessageEvent(chatname, username, t)}
}

func (e *SystemMessageChatLeaveEvent) String() string {
	return fmt.Sprintf("sysMsg: chat leave: %v", e.SystemMessageEvent)
}
