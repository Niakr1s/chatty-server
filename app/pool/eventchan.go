package pool

import "github.com/niakr1s/chatty-server/app/pool/events"

// EventChan is used to determine if event can pass into channel
// custome FilterPass can be assigned, which returns true if event is passable
type EventChan struct {
	Ch chan events.Event

	// should return true if event should pass into channel
	Filter FilterPass
}

// NewEventChan ...
func NewEventChan() *EventChan {
	return &EventChan{Ch: make(chan events.Event, 10), Filter: FilterPassAlways}
}

// WithFilter ...
func (ec *EventChan) WithFilter(filter FilterPass) *EventChan {
	ec.Filter = filter
	return ec
}
