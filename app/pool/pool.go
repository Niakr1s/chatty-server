package pool

import (
	"server2/app/er"
	"server2/app/pool/events"
	"server2/app/store"
	"sync"
)

// Pool ...
type Pool struct {
	sync.Mutex

	store *store.Store

	// events.Event inputs here
	inputCh chan events.Event

	userCh       map[string]*events.EventChan
	userChFilter func(username string) events.FilterPass

	innerCh []*events.EventChan
}

// NewPool ...
func NewPool(s *store.Store) *Pool {
	return &Pool{
		store: s,

		inputCh: make(chan events.Event, 10),

		userCh:       make(map[string]*events.EventChan),
		userChFilter: func(username string) events.FilterPass { return FilterPassIfUserInChat(s, username) },

		innerCh: make([]*events.EventChan, 0),
	}
}

// WithUserChFilter ...
func (p *Pool) WithUserChFilter(f func(username string) events.FilterPass) *Pool {
	p.userChFilter = f
	return p
}

// GetInputChan ...
func (p *Pool) GetInputChan() chan<- events.Event {
	return p.inputCh
}

// GetUserChan gets chan for user
// chans are created with default filter events.FilterPassIfUserInChat
func (p *Pool) GetUserChan(username string) <-chan events.Event {
	p.Lock()
	defer p.Unlock()

	return p.getUserChan(username)
}

// GetUserChan gets chan for user
// chans are created with default filter events.FilterPassIfUserInChat
func (p *Pool) getUserChan(username string) <-chan events.Event {
	if _, ok := p.userCh[username]; !ok {
		p.createUserChan(username)
	}

	return p.userCh[username].Ch
}

func (p *Pool) createUserChan(username string) {
	p.userCh[username] = events.NewEventChan().
		WithFilter(p.userChFilter(username))
}

func (p *Pool) removeUserChan(username string) {
	if u, ok := p.userCh[username]; ok {
		close(u.Ch)
		delete(p.userCh, username)
	}
}

// CreateChanNoFilter creates chan with no filter
func (p *Pool) CreateChanNoFilter() <-chan events.Event {
	return p.CreateChan(events.FilterPassAlways)
}

// CreateChan creates chan with filter
func (p *Pool) CreateChan(filter events.FilterPass) <-chan events.Event {
	p.Lock()
	defer p.Unlock()

	return p.createChan(filter)
}

// CreateChan creates chan with filter
func (p *Pool) createChan(filter events.FilterPass) <-chan events.Event {
	ec := events.NewEventChan().WithFilter(filter)

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

func (p *Pool) beforeSending(event events.Event) {
	p.processLogoutEvent(event)
	p.processLoginEvent(event)
}

func (p *Pool) processLogoutEvent(event events.Event) {
	if logout, ok := event.(*events.LogoutEvent); ok {
		if _, err := logout.InChat(); err == er.ErrGlobalEvent {
			p.removeUserChan(logout.Username)
		}
	}
}

func (p *Pool) processLoginEvent(event events.Event) {
	if login, ok := event.(*events.LoginEvent); ok {
		if _, err := login.InChat(); err == er.ErrGlobalEvent {
			p.removeUserChan(login.Username) // to renew channel
			p.createUserChan(login.Username)
		}
	}
}

func (p *Pool) sendInUserChans(event events.Event) {
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

func (p *Pool) sendInInnerChans(event events.Event) {
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
