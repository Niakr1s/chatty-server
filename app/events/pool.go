package events

import (
	"server2/app/store"
	"sync"
)

// Pool ...
type Pool struct {
	sync.Mutex

	store *store.Store

	// Event inputs here
	inputCh chan Event

	userCh       map[string]*EventChan
	userChFilter func(username string) FilterPass

	innerCh []*EventChan
}

// NewPool ...
func NewPool(s *store.Store) *Pool {
	return &Pool{
		store: s,

		inputCh: make(chan Event, 10),

		userCh:       make(map[string]*EventChan),
		userChFilter: func(username string) FilterPass { return FilterPassIfUserInChat(s, username) },

		innerCh: make([]*EventChan, 0),
	}
}

// WithUserChFilter ...
func (p *Pool) WithUserChFilter(f func(username string) FilterPass) *Pool {
	p.userChFilter = f
	return p
}

// GetInputChan ...
func (p *Pool) GetInputChan() chan<- Event {
	return p.inputCh
}

// GetUserChan gets chan for user
// chans are created with default filter FilterPassIfUserInChat
func (p *Pool) GetUserChan(username string) <-chan Event {
	p.Lock()
	defer p.Unlock()

	return p.getUserChan(username)
}

// GetUserChan gets chan for user
// chans are created with default filter FilterPassIfUserInChat
func (p *Pool) getUserChan(username string) <-chan Event {
	if _, ok := p.userCh[username]; !ok {
		p.createUserChan(username)
	}

	return p.userCh[username].Ch
}

func (p *Pool) createUserChan(username string) {
	p.userCh[username] = NewEventChan().
		WithFilter(p.userChFilter(username))
}

func (p *Pool) removeUserChan(username string) {
	if u, ok := p.userCh[username]; ok {
		close(u.Ch)
		delete(p.userCh, username)
	}
}

// CreateChanNoFilter creates chan with no filter
func (p *Pool) CreateChanNoFilter() <-chan Event {
	return p.CreateChan(FilterPassAlways)
}

// CreateChan creates chan with filter
func (p *Pool) CreateChan(filter FilterPass) <-chan Event {
	p.Lock()
	defer p.Unlock()

	return p.createChan(filter)
}

// CreateChan creates chan with filter
func (p *Pool) createChan(filter FilterPass) <-chan Event {
	ec := NewEventChan().WithFilter(filter)

	p.innerCh = append(p.innerCh, ec)

	return ec.Ch
}

// Run ...
func (p *Pool) Run() {
	go func() {
		for {
			event := <-p.inputCh

			p.Lock()

			p.beforeSending(event)

			p.sendInUserChans(event)
			p.sendInInnerChans(event)

			p.Unlock()
		}
	}()
}

func (p *Pool) beforeSending(event Event) {
	p.processLogoutEvent(event)
	p.processLoginEvent(event)
}

func (p *Pool) processLogoutEvent(event Event) {
	if logout, ok := event.(*LogoutEvent); ok {
		if _, err := logout.InChat(); err == ErrGlobal {
			p.removeUserChan(logout.Username)
		}
	}
}

func (p *Pool) processLoginEvent(event Event) {
	if login, ok := event.(*LoginEvent); ok {
		if _, err := login.InChat(); err == ErrGlobal {
			p.removeUserChan(login.Username) // to renew channel
			p.createUserChan(login.Username)
		}
	}
}

func (p *Pool) sendInUserChans(event Event) {
	for k, eventCh := range p.userCh {
		if eventCh.Filter(event) {
			select {
			case eventCh.Ch <- event:
			default:
				// if channel is full - close and delete it, it's safe for map
				p.removeUserChan(k)
			}
		}
	}
}

func (p *Pool) sendInInnerChans(event Event) {
	n := len(p.innerCh)
	for i := 0; i < n; i++ {
		eventCh := p.innerCh[i]
		if eventCh.Filter(event) {
			select {
			case eventCh.Ch <- event:
			default:
				// if channel is full - close it and set it to nil
				close(eventCh.Ch)
				eventCh = nil

				// swapping with last element
				p.innerCh[i], p.innerCh[n-1] = p.innerCh[n-1], p.innerCh[i]
				i--
				n--
			}
		}
	}
	p.innerCh = p.innerCh[:n]
}
