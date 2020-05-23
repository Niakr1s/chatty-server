package er

import "errors"

// Errors
var (
	ErrPasswordIsEmpty     = errors.New("password is empty")
	ErrPasswordHashIsEmpty = errors.New("password hash is empty")

	ErrHashMismatch = errors.New("hash mismatch")
)
