package events

import "fmt"

// Event ...
type Event interface {
	fmt.Stringer
	// should return chat this event occured in
	InChat() (string, error)
}
