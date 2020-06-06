package db

import (
	"github.com/niakr1s/chatty-server/app/models"
)

// MessageDB ...
type MessageDB interface {
	// should update message ID and time
	Post(*models.Message) error

	// should return empty slice even on error
	GetLastNMessages(chatname string, n int) ([]*models.Message, error)
}
