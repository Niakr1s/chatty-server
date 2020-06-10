package command

import (
	"errors"
)

// errors
var (
	ErrBadInput      = errors.New("Empty text")
	ErrNoSuchCommand = errors.New("no such command")
)

// Command used by bot to answer user.
// Bot can response to commands of type "/help", or maybe "/calc 2+2".
type Command interface {
	Answer() (string, error)
}

// CommandFunc is adapter for Command
type CommandFunc func() (string, error)

// Answer ...
func (cf CommandFunc) Answer() (string, error) {
	return cf()
}
