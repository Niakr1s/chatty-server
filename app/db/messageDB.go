package db

import (
	"sync"

	"github.com/niakr1s/chatty-server/app/models"
)

// MessageDB ...
type MessageDB interface {
	sync.Locker

	// should update message ID
	Post(*models.Message) error

	GetLastNMessages(chatname string, n int) ([]*models.Message, error)
}
