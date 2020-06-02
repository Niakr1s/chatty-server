package er

import (
	"errors"
	"fmt"

	"github.com/niakr1s/chatty-server/app/constants"
)

// Chat DB errors
var (
	ErrChatAlreadyExists = errors.New("chat already exists")
	ErrNoSuchChat        = errors.New("no such chat")
)

// User DB errors
var (
	ErrUserNotFound          = errors.New("user doesn't exist")
	ErrUserAlreadyRegistered = errors.New("user already registered")
	ErrCannotStoreUser       = errors.New("cannot store user")
	ErrCannotUpdateUser      = errors.New("cannot update user")
)

// Logged DB errors
var (
	ErrAlreadyLogged = errors.New("user is already logged in")
	ErrNotLogged     = errors.New("user not logged in")
)

// Events errors
var (
	ErrGlobalEvent = errors.New("global event")
)

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

// Server errors
var (
	ErrCannotParseData    = errors.New("cannot parse data")
	ErrUnathorized        = errors.New("unauthorized")
	ErrNoUsername         = errors.New("no username")
	ErrNoActivationToken  = errors.New("no activation token")
	ErrBadActivationToken = errors.New("bad activation token")
)

// Email errors
var (
	ErrSendEmail       = errors.New("couldn't send email")
	ErrUnverifiedEmail = errors.New("email is not verified")
)

// Session errors
var (
	ErrSession = errors.New("couldn't get session")
)

// Type errors
var (
	ErrConvertType = errors.New("couldn't convert type")
)

// Env errors
var (
	ErrEnvEmptySendGridAPIKey = fmt.Errorf("empty %s", constants.EnvSendGridAPIKey)
)
