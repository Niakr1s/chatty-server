package logged

import "errors"

// DB errors
var (
	ErrAlreadyLogged = errors.New("user is already logged in")
	ErrNotLogged     = errors.New("user not logged in")
)
