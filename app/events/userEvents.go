package events

import (
	"fmt"
	"time"

	"github.com/niakr1s/chatty-server/app/models"
)

// UserEvent ...
type UserEvent struct {
	models.User
	models.Chat
	Time models.UnixTime
}

// NewUserEvent ...
func NewUserEvent(username string, chatname string, time time.Time) *UserEvent {
	return &UserEvent{User: models.NewUser(username), Chat: models.NewChat(chatname), Time: models.UnixTime(time)}
}

func (ue *UserEvent) String() string {
	return fmt.Sprintf("user %s, chat: %s, time: %v", ue.UserName, ue.ChatName, ue.Time)
}

// LoginEvent used to emit when user logins chat globally (ie to db.LoggedDB)
type LoginEvent struct {
	*UserEvent
}

// NewLoginEvent ...
func NewLoginEvent(username string, chatname string, time time.Time) *LoginEvent {
	return &LoginEvent{NewUserEvent(username, chatname, time)}
}

func (le *LoginEvent) String() string {
	return fmt.Sprintf("login: %v", le.UserEvent)
}

// LogoutEvent used to emit when user logouts globally (ie from db.LoggedDB)
type LogoutEvent struct {
	*UserEvent
}

// NewLogoutEvent ...
func NewLogoutEvent(username string, chatname string, time time.Time) *LogoutEvent {
	return &LogoutEvent{NewUserEvent(username, chatname, time)}
}

func (le *LogoutEvent) String() string {
	return fmt.Sprintf("logout: %v", le.UserEvent)
}

// ChatJoinEvent used to emit when user joins chat (ie to db.Chat)
type ChatJoinEvent struct {
	*UserEvent
}

// NewChatJoinEvent ...
func NewChatJoinEvent(username string, chatname string, time time.Time) *ChatJoinEvent {
	return &ChatJoinEvent{NewUserEvent(username, chatname, time)}
}

func (ce *ChatJoinEvent) String() string {
	return fmt.Sprintf("join chat: %v", ce.UserEvent)
}

// ChatLeaveEvent used to emit when user leaves chat (ie from db.Chat)
type ChatLeaveEvent struct {
	*UserEvent
}

// NewChatLeaveEvent ...
func NewChatLeaveEvent(username string, chatname string, time time.Time) *ChatLeaveEvent {
	return &ChatLeaveEvent{NewUserEvent(username, chatname, time)}
}

func (ce *ChatLeaveEvent) String() string {
	return fmt.Sprintf("leave chat: %v", ce.UserEvent)
}
