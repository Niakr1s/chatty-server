package bot

import (
	"context"
	"time"

	"github.com/niakr1s/chatty-server/app/bot/internal/command"
	"github.com/niakr1s/chatty-server/app/client"
	"github.com/niakr1s/chatty-server/app/events"
	log "github.com/sirupsen/logrus"
)

const (
	contentJSON = "application/json"
)

// Bot ...
type Bot struct {
	*client.Client
	ctx context.Context
}

// New constructs a Bot
func New(ctx context.Context, botname, password, url string) (*Bot, error) {
	client, err := client.New(ctx, botname, password, url)
	if err != nil {
		return nil, err
	}
	return &Bot{Client: client, ctx: ctx}, nil
}

// Run ...
func (b *Bot) Run() error {
	// start sending keep alive, we'll do it forever
	b.StartSendingKeepAlive()

	// start our loop
	for b.loop() {
	}

	return nil
}

// loop returns true if need to run again, otherwise false.
func (b *Bot) loop() bool {
	// connecting...
	for b.Connect() != nil {
		log.Infof("couldn't reconnect to %s, trying in 10secs...", b.URL)
		select {
		case <-b.ctx.Done():
			return false
		case <-time.After(time.Second * 10):
		}
	}

	chats, err := b.GetChats()
	if err != nil {
		return true
	}
	b.JoinChats(chats...)
	eventsCh := b.StartListen()
	for {
		select {
		case <-b.ctx.Done():
			return false
		case e, ok := <-eventsCh:
			// channel has been closed, go reconnect
			if !ok {
				return true
			}
			log.Tracef("got event: %v", e)
			switch e := e.(type) {
			case *events.ChatJoinEvent:
			case *events.MessageEvent:
				cmd, err := command.ParseCommand(b.Username, e.Message)
				if err != nil {
					continue
				}
				answer, err := cmd.Answer()
				if err != nil {
					continue
				}
				b.PostMessage(e.ChatName, answer)
			case *events.ChatCreatedEvent:
				go func() { b.JoinChat(e.ChatName) }()
			default:
			}
		}
	}
}
