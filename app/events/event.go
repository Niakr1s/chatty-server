package events

import "fmt"

// Event ...
type Event interface {
	fmt.Stringer
	// should return chat this event occured in
	// if it's global, should return "", er.ErrGlobalEvent
	InChat() (string, error)
}
