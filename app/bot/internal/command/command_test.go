package command

import (
	"fmt"
	"testing"

	"github.com/niakr1s/chatty-server/app/models"
	"github.com/stretchr/testify/assert"
)

func TestParseCommand(t *testing.T) {
	const botName = "Bot"

	tests := []struct {
		name    string
		message *models.Message
		wantErr bool
	}{
		{"valid user, valid text with comma", models.NewMessage("no bot", fmt.Sprintf("%s, /help", botName), ""), false},
		{"valid user, valid text without comma", models.NewMessage("no bot", fmt.Sprintf("%s /help", botName), ""), false},
		{"valid user, asking bot with invalid format", models.NewMessage("no bot", fmt.Sprintf("%s help", botName), ""), true},
		{"valid user, asking no bot", models.NewMessage("no bot", fmt.Sprintf("%s, /help", "nobot"), ""), true},
		{"Bot asking Bot", models.NewMessage(botName, fmt.Sprintf("%s, /help", botName), ""), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseCommandForBot(botName, tt.message)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
