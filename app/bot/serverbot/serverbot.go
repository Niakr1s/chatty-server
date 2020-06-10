package serverbot

import (
	"github.com/niakr1s/chatty-server/app/bot/internal/command"
	"github.com/niakr1s/chatty-server/app/db"
	"github.com/niakr1s/chatty-server/app/events"
	"github.com/niakr1s/chatty-server/app/models"
	log "github.com/sirupsen/logrus"
)

// Bot ...
type Bot struct {
	inputCh <-chan events.Event

	closeCh chan struct{}
	doneCh  chan struct{}

	messageDB db.MessageDB
}

// New constructs a Bot
func New(inputCh <-chan events.Event, messageDB db.MessageDB) *Bot {
	return &Bot{inputCh: inputCh, closeCh: make(chan struct{}), doneCh: make(chan struct{}), messageDB: messageDB}
}

// Run ...
func (b *Bot) Run() {
	// start our loop
	go func() {
		for b.loop() {
		}
		close(b.doneCh)
	}()
}

// Shutdown ...
func (b *Bot) Shutdown() <-chan struct{} {
	close(b.closeCh)
	return b.doneCh
}

// loop returns true if need to run again, otherwise false.
func (b *Bot) loop() bool {
	for {
		select {
		case <-b.closeCh:
			return false

		case e, ok := <-b.inputCh:
			// channel has been closed, go reconnect
			if !ok {
				return true
			}
			log.Tracef("got event: %v", e)
			switch e := e.(type) {
			case *events.ChatJoinEvent:
				// maybe greet user
			case *events.MessageEvent:
				cmd, err := command.ParseCommand(e.Message)
				if err != nil {
					continue
				}
				answer, err := cmd.Answer()
				if err != nil {
					continue
				}
				if err := b.messageDB.Post(models.NewMessage("Bot", answer, e.ChatName)); err != nil {
					log.Errorf("bot: couldn't post message: %v", err)
				}
			default:
			}
		}
	}
}
