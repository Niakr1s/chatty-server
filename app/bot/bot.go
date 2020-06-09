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

// HelloBot ...
type HelloBot struct {
	client *http.Client

	ctx context.Context
	url string

	botName string
}

// NewHelloBot ...
func NewHelloBot() (*HelloBot, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	return &HelloBot{client: &http.Client{Jar: jar}}, nil
}

// Connect ...
func (b *HelloBot) Connect(botname, password, url string) error {
	b.botName = botname
	b.url = url
	w, err := b.client.Post(url+constants.RouteApi+constants.RouteAuthorize,
		contentJSON, strings.NewReader(fmt.Sprintf(`{"user": "%s", "password": "%s"}`, botname, password)))
	if err != nil {
		return err
	}
	w.Body.Close()
	return nil
}

// Run ...
func (b *HelloBot) Run(ctx context.Context) error {
	b.ctx = ctx
	chats, err := b.getChats()
	if err != nil {
		return err
	}
	b.joinChats(chats...)
	b.startSendingKeepAlive()
	eventsCh := b.startListen()
	for {
		select {
		case <-ctx.Done():
			return nil
		case e := <-eventsCh:
			log.Tracef("got event: %v", e)
			switch e := e.(type) {
			case *events.ChatJoinEvent:
			case *events.MessageEvent:
				cmd, err := command.ParseCommand(b.botName, e.Message)
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

func (b *HelloBot) postHelpMessage(chat string) {
	b.postMessage(chat, fmt.Sprintf(`Usage info: post message with "%s, /command" to invoke command.
Available commands:
	help: prints this message`, b.botName))
}

func (b *HelloBot) greetUser(chat, user string) {
	b.postMessage(chat, fmt.Sprintf("Hello, %s, how are you?", user))
}

func (b *HelloBot) startListen() <-chan events.Event {
	ch := make(chan events.Event)
	go func() {
		pollURL := b.url + constants.RouteApi + constants.RouteLoggedOnly + constants.RoutePoll
		for {
			select {
			case <-b.ctx.Done():
				return
			default:
				w, err := b.client.Get(pollURL)
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

func (b *HelloBot) postMessage(chat, text string) {
	log.Tracef("posting message in chat %s: %s", chat, text)
	go func() {
		postMsgURL := b.url + constants.RouteApi + constants.RouteLoggedOnly + constants.RoutePostMessage
		if _, err := b.client.Post(postMsgURL, contentJSON,
			strings.NewReader(fmt.Sprintf(`{"user": "%s", "chat": "%s", "text": "%s"}`,
				b.botName, chat, text))); err != nil {
			log.Errorf("postMessage: %v", err)
		}
	}()
}

func (b *HelloBot) getChats() ([]string, error) {
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

func (b *HelloBot) joinChats(chatnames ...string) {
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

func (b *HelloBot) joinChat(chat string) error {
	joinChatURL := b.url + constants.RouteApi + constants.RouteLoggedOnly + constants.RouteJoinChat
	_, err := b.client.Post(joinChatURL, contentJSON, strings.NewReader(fmt.Sprintf(`{"chat": "%s"}`, chat)))
	if err != nil {
		return ErrJoinChat
	}
	return nil
}

func (b *HelloBot) startSendingKeepAlive() {
	go func() {
		for {
			<-time.After(time.Second * 10)
			if err := b.sendKeepAlive(); err != nil {
				log.Errorf("send alive: %v", err)
			}
		}
	}()
}

func (b *HelloBot) sendKeepAlive() error {
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
