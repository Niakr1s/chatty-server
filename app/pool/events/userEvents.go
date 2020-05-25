package events

import (
	"fmt"
	"time"

	"github.com/niakr1s/chatty-server/app/er"
)

// UserEvent ...
type UserEvent struct {
	Username string
	Chatname string
	Time     time.Time
}

// NewUserEvent ...
func NewUserEvent(username string, chatname string, time time.Time) *UserEvent {
	return &UserEvent{Username: username, Chatname: chatname, Time: time}
}

// InChat ...
func (ue *UserEvent) InChat() (string, error) {
	if ue.Chatname == "" {
		return "", er.ErrGlobalEvent
	}
	return ue.Chatname, nil
}

func (ue *UserEvent) String() string {
	return fmt.Sprintf("user %s, chat: %s, time: %v", ue.Username, ue.Chatname, ue.Time)
}

// LoginEvent ...
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

// LogoutEvent ...
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

// ChatJoinEvent ...
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

// ChatLeaveEvent ...
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
