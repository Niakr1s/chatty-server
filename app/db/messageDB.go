package db

import (
	"server2/app/models"
	"sync"
)

// MessageDB ...
type MessageDB interface {
	sync.Locker

	// should update message ID
	Post(*models.Message) error

	GetLastNMessages(chatname string, n int) ([]*models.Message, error)
}
