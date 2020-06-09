package command

import (
	"errors"
	"strings"

	"github.com/niakr1s/chatty-server/app/models"
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

// ParseCommand gots input of type "Bot, /help smth" (comma after bot name can be omitted) and returns a command
func ParseCommand(botname string, msg *models.Message) (Command, error) {
	if len(msg.Text) == 0 || msg.UserName == botname {
		return nil, ErrBadInput
	}
	splitted := strings.SplitN(msg.Text, " ", 2)
	if len(splitted) < 2 || !(splitted[0] == botname+"," || splitted[0] == botname) {
		return nil, ErrBadInput
	}

	cmd := splitted[1]
	// arg := ""
	// if len(splitted) > 2 {
	// 	arg = splitted[2]
	// }

	switch cmd {
	case "/help":
		return HelpCommand(botname), nil
	case "/anecdot", "/anekdot":
		return AnecdotCommand(), nil
	default:
		return nil, ErrNoSuchCommand
	}
}
