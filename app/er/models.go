package er

import "errors"

// Chat errors
var (
	ErrAlreadyInChat = errors.New("user is already logged in chat")
	ErrNotInChat     = errors.New("user not logged in such chat")
)

// User errors
var (
	ErrPasswordIsEmpty     = errors.New("password is empty")
	ErrPasswordHashIsEmpty = errors.New("password hash is empty")

	ErrHashMismatch = errors.New("hash mismatch")

	ErrUserNameIsEmpty = errors.New("username is empty")
)

// Message errors
var (
	ErrTooOld = errors.New("too old")
)
