package events

import (
	"fmt"
	"reflect"
)

// Event ...
type Event interface {
	fmt.Stringer
	// should return chat this event occured in
	// if it's global, should return "", er.ErrGlobalEvent
	InChat() (string, error)
}

// EventWithType represents event with reflected Type
// Very convinient to json marshalling for example
type EventWithType struct {
	Event `json:"event"`

	Type string `json:"type"`
}

// NewEventWithType constructs EventWithType from Event
func NewEventWithType(e Event) EventWithType {
	return EventWithType{Event: e, Type: reflect.TypeOf(e).Elem().Name()}
}
