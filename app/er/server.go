package er

import "errors"

// Server errors
var (
	ErrCannotParseData = errors.New("cannot parse data")
	ErrUnathorized     = errors.New("unauthorized")
)
