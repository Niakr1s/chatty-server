package events

import "time"

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
		return "", ErrGlobal
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
