package models

import (
	"fmt"
	"time"
)

// UnixTime is wrapper fot time.Time that marshals into epoch unix time
type UnixTime time.Time

// MarshalJSON impl
func (t UnixTime) MarshalJSON() ([]byte, error) {
	str := fmt.Sprintf("%d", time.Time(t).Unix())
	return []byte(str), nil
}
