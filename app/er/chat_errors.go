package er

import "errors"

// Chat errors
var (
	ErrAlreadyInChat = errors.New("user is already logged in chat")
	ErrNotInChat     = errors.New("user not logged in such chat")

	ErrChatAlreadyExists = errors.New("chat already exists")
	ErrNoSuchChat        = errors.New("no such chat")
)
