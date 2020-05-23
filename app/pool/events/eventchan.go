package events

// EventChan ...
type EventChan struct {
	Ch chan Event

	// should return true if event should pass into channel
	Filter FilterPass
}

// NewEventChan ...
func NewEventChan() *EventChan {
	return &EventChan{Ch: make(chan Event, 10), Filter: FilterPassAlways}
}

// WithFilter ...
func (ec *EventChan) WithFilter(filter FilterPass) *EventChan {
	ec.Filter = filter
	return ec
}
