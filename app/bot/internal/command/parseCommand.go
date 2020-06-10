package command

import (
	"strings"

	"github.com/niakr1s/chatty-server/app/models"
)

// ParseCommandForBot gots input of type "Bot, /help smth" (comma after bot name can be omitted) and returns a command
func ParseCommandForBot(botname string, msg *models.Message) (Command, error) {
	if len(msg.Text) == 0 || msg.UserName == botname {
		return nil, ErrBadInput
	}
	splitted := strings.SplitN(msg.Text, " ", 2)
	if len(splitted) < 2 || !(splitted[0] == botname+"," || splitted[0] == botname) {
		return nil, ErrBadInput
	}

	cmd := splitted[1]
	arg := ""
	return getCommand(cmd, arg)
}

// ParseCommand gots input of type "/help smth" and returns a command
func ParseCommand(msg *models.Message) (Command, error) {
	splitted := strings.SplitN(msg.Text, " ", 1)
	if len(splitted) == 0 {
		return nil, ErrBadInput
	}

	cmd := splitted[1]
	arg := ""
	return getCommand(cmd, arg)
}

func getCommand(cmd, arg string) (Command, error) {
	switch cmd {
	case "/help":
		return HelpCommand(), nil
	case "/anecdot", "/anekdot":
		return AnecdotCommand(), nil
	default:
		return nil, ErrNoSuchCommand
	}
}
