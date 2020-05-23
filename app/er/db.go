package er

import "errors"

// Chat DB errors
var (
	ErrChatAlreadyExists = errors.New("chat already exists")
	ErrNoSuchChat        = errors.New("no such chat")
)

// User DB errors
var (
	ErrUserNotFound = errors.New("user doesn't exist")
)

// Logged DB errors
var (
	ErrAlreadyLogged = errors.New("user is already logged in")
	ErrNotLogged     = errors.New("user not logged in")
)
