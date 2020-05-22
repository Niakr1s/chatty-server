package logged

import "errors"

// DB errors
var (
	ErrAlreadyLogged = errors.New("user is already logged in")
	ErrNotLogged     = errors.New("user not logged in")
)

// User errors
var (
	ErrAlreadyInChat = errors.New("user is already logged in chat")
	ErrNotInChat     = errors.New("user not logged in such chat")
)
