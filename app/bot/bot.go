package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"sync"
	"time"

	"github.com/niakr1s/chatty-server/app/bot/internal/command"
	"github.com/niakr1s/chatty-server/app/constants"
	"github.com/niakr1s/chatty-server/app/db"
	"github.com/niakr1s/chatty-server/app/events"
	log "github.com/sirupsen/logrus"
)

const (
	contentJSON = "application/json"
)

// Bot ...
type Bot struct {
	client *http.Client

	ctx context.Context

	botname  string
	password string
	url      string
}

// New ...
func New(ctx context.Context, botname, password, url string) (*Bot, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	return &Bot{client: &http.Client{Jar: jar}, botname: botname, password: password, url: url, ctx: ctx}, nil
}

// connect ...
func (b *Bot) connect() error {
	w, err := b.client.Post(b.url+constants.RouteApi+constants.RouteAuthorize,
		contentJSON, strings.NewReader(fmt.Sprintf(`{"user": "%s", "password": "%s"}`, b.botname, b.password)))
	if err != nil {
		return err
	}
	defer w.Body.Close()
	return nil
}

// Run ...
func (b *Bot) Run() error {
	// start sending keep alive, we'll do it forever
	b.startSendingKeepAlive()

	// start our loop
	for b.loop() {
	}

	return nil
}

// loop returns true if need to run again, otherwise false.
func (b *Bot) loop() bool {
	// connecting...
	for b.connect() != nil {
		log.Infof("couldn't reconnect, trying in 10secs...")
		<-time.After(time.Second * 10)
	}

	chats, err := b.getChats()
	if err != nil {
		return true
	}
	b.joinChats(chats...)
	eventsCh := b.startListen()
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
				cmd, err := command.ParseCommand(b.botname, e.Message)
				if err != nil {
					continue
				}
				answer, err := cmd.Answer()
				if err != nil {
					continue
				}
				b.postMessage(e.ChatName, answer)
			default:
			}
		}
	}
}

func (b *Bot) postHelpMessage(chat string) {
	b.postMessage(chat, fmt.Sprintf(`Usage info: post message with "%s, /command" to invoke command.
Available commands:
	help: prints this message`, b.botname))
}

func (b *Bot) greetUser(chat, user string) {
	b.postMessage(chat, fmt.Sprintf("Hello, %s, how are you?", user))
}

func (b *Bot) startListen() <-chan events.Event {
	ch := make(chan events.Event)
	go func() {
		pollURL := b.url + constants.RouteApi + constants.RouteLoggedOnly + constants.RoutePoll
		for {
			select {
			case <-b.ctx.Done():
				return
			default:
				w, err := b.client.Get(pollURL)
				if err != nil {
					close(ch)
					return
				}
				e, err := parseEvent(w.Body)
				if err != nil {
					continue
				}
				ch <- e
			}
		}
	}()
	return ch
}

func (b *Bot) postMessage(chat, text string) {
	log.Tracef("posting message in chat %s: %s", chat, text)
	go func() {
		postMsgURL := b.url + constants.RouteApi + constants.RouteLoggedOnly + constants.RoutePostMessage
		if _, err := b.client.Post(postMsgURL, contentJSON,
			strings.NewReader(fmt.Sprintf(`{"user": "%s", "chat": "%s", "text": "%s"}`,
				b.botname, chat, text))); err != nil {
			log.Errorf("postMessage: %v", err)
		}
	}()
}

func (b *Bot) getChats() ([]string, error) {
	getChatsURL := b.url + constants.RouteApi + constants.RouteLoggedOnly + constants.RouteGetChats

	w, err := b.client.Get(getChatsURL)
	if err != nil {
		return nil, err
	}
	defer w.Body.Close()

	chatReports := []db.ChatReport{}
	if err := json.NewDecoder(w.Body).Decode(&chatReports); err != nil {
		return nil, err
	}
	res := make([]string, 0, len(chatReports))
	for _, cr := range chatReports {
		res = append(res, cr.ChatName)
	}
	log.Tracef("got %d chats", len(res))
	return res, nil
}

func (b *Bot) joinChats(chatnames ...string) {
	wg := sync.WaitGroup{}
	wg.Add(len(chatnames))
	for _, cn := range chatnames {
		go func(chatname string) {
			defer wg.Done()
			if err := b.joinChat(chatname); err != nil {
				log.Error(err)
				return
			}
		}(cn)
	}
	wg.Wait()
}

func (b *Bot) joinChat(chat string) error {
	joinChatURL := b.url + constants.RouteApi + constants.RouteLoggedOnly + constants.RouteJoinChat
	_, err := b.client.Post(joinChatURL, contentJSON, strings.NewReader(fmt.Sprintf(`{"chat": "%s"}`, chat)))
	if err != nil {
		return ErrJoinChat
	}
	return nil
}

func (b *Bot) startSendingKeepAlive() {
	go func() {
		for {
			<-time.After(time.Second * 10)
			b.sendKeepAlive()
		}
	}()
}

func (b *Bot) sendKeepAlive() error {
	keepAliveURL := b.url + constants.RouteApi + constants.RouteLoggedOnly + constants.RouteKeepAlive
	r, err := http.NewRequestWithContext(b.ctx, http.MethodPut, keepAliveURL, nil)
	if err != nil {
		return err
	}
	if _, err := b.client.Do(r); err != nil {
		return err
	}
	return nil
}
