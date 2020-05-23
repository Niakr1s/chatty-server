package events

// Event ...
type Event interface {
	// should return chat this event occured in
	InChat() (string, error)
}
