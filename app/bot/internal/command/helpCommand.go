package command

// HelpCommand for string "/help"
func HelpCommand(botname string) CommandFunc {
	return func() (string, error) {
		return `Available commands: /help`, nil
	}
}
