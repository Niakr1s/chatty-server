package pool

import (
	"fmt"
	"server2/app/pool/events"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockEvent struct{}

const chat = "chat"
const username = "user"

func (e *mockEvent) InChat() (string, error) {
	return chat, nil
}

func NewMockPool(t *testing.T, inner int, user int) *Pool {
	t.Helper()

	p := NewPool(nil).WithUserChFilter(func(username string) events.FilterPass { return events.FilterPassAlways })

	for i := 0; i < inner; i++ {
		p.CreateChanNoFilter()
	}

	for i := 0; i < user; i++ {
		p.GetUserChan(fmt.Sprintf("%s%d", username, i))
	}

	return p
}

func TestPool_GetUserChan(t *testing.T) {
	p := NewMockPool(t, 0, 0)
	assert.Len(t, p.userCh, 0)

	p.GetUserChan(username)
	assert.Len(t, p.userCh, 1)

	p.GetUserChan(username)
	assert.Len(t, p.userCh, 1)

	p.GetUserChan(username + "1")
	assert.Len(t, p.userCh, 2)
}

func TestPool_CreateChan(t *testing.T) {
	p := NewMockPool(t, 0, 0)
	assert.Len(t, p.innerCh, 0)

	p.CreateChanNoFilter()
	assert.Len(t, p.innerCh, 1)

	p.CreateChanNoFilter()
	assert.Len(t, p.innerCh, 2)
}

func TestPool_send(t *testing.T) {
	for i := 1; i != 3; i++ {
		for j := 1; j != 3; j++ {
			i := i
			j := j
			t.Run("succesfully created chans", func(t *testing.T) {
				p := NewMockPool(t, i, j)

				assert.NotEmpty(t, p.innerCh)
				assert.NotEmpty(t, p.userCh)
			})
		}
	}

	for i := 0; i != 3; i++ {
		for j := 0; j != 3; j++ {
			i := i
			j := j
			t.Run(fmt.Sprintf("%d inner and %d user chans were emptied", i, j), func(t *testing.T) {
				p := NewMockPool(t, i, j)

				for i := 0; i < 100; i++ {
					p.sendInInnerChans(&mockEvent{})
					p.sendInUserChans(&mockEvent{})
				}
				assert.Empty(t, p.innerCh)
			})
		}
	}

}

func TestPool_ProcessLogoutEvent(t *testing.T) {
	p := NewMockPool(t, 0, 0)

	p.GetUserChan(username)

	assert.Len(t, p.userCh, 1)

	logoutEvent := &events.LogoutEvent{UserEvent: &events.UserEvent{Username: username, Chatname: ""}}

	p.processLogoutEvent(logoutEvent)

	assert.Len(t, p.userCh, 0)
}

func TestPool_Run(t *testing.T) {
	p := NewMockPool(t, 0, 0)
	p.Run()

	inputCh := p.GetInputChan()

	userCh := p.GetUserChan(username)
	innerCh := p.CreateChanNoFilter()

	for i := 0; i < 3; i++ {
		event := &mockEvent{}
		inputCh <- event
	}

	for i := 0; i < 3; i++ {
		assert.NotNil(t, <-userCh)
		assert.NotNil(t, <-innerCh)
	}

	select {
	case <-userCh:
		assert.Fail(t, "non-empty userCh")
	case <-innerCh:
		assert.Fail(t, "non-empty innerCh")
	default:
	}

}
