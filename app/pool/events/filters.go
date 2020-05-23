package events

// FilterPass ...
type FilterPass func(e Event) bool

// FilterPassAlways ...
func FilterPassAlways(e Event) bool {
	return true
}
