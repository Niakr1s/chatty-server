package events

import (
	"server2/app/er"
	"time"
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

// LoginEvent ...
type LoginEvent struct {
	*UserEvent
}

// NewLoginEvent ...
func NewLoginEvent(username string, chatname string, time time.Time) *LoginEvent {
	return &LoginEvent{NewUserEvent(username, chatname, time)}
}

// LogoutEvent ...
type LogoutEvent struct {
	*UserEvent
}

// NewLogoutEvent ...
func NewLogoutEvent(username string, chatname string, time time.Time) *LogoutEvent {
	return &LogoutEvent{NewUserEvent(username, chatname, time)}
}

// ChatJoinEvent ...
type ChatJoinEvent struct {
	*UserEvent
}

// NewChatJoinEvent ...
func NewChatJoinEvent(username string, chatname string, time time.Time) *ChatJoinEvent {
	return &ChatJoinEvent{NewUserEvent(username, chatname, time)}
}

// ChatLeaveEvent ...
type ChatLeaveEvent struct {
	*UserEvent
}

// NewChatLeaveEvent ...
func NewChatLeaveEvent(username string, chatname string, time time.Time) *ChatLeaveEvent {
	return &ChatLeaveEvent{NewUserEvent(username, chatname, time)}
}
