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

// ToSQLTimeStamp ...
func (t UnixTime) ToSQLTimeStamp() string {
	return time.Time(t).Format("2006-01-02 15:04:05")
}
