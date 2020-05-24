package er

import "errors"

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
	ErrSendEmail = errors.New("couldn't send email")
)
