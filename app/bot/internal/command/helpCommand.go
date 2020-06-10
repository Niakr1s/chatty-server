package command

// HelpCommand for string "/help"
func HelpCommand() CommandFunc {
	return func() (string, error) {
		return `Available commands: /help, /anecdot (/anekdot)`, nil
	}
}
