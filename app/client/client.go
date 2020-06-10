package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"sync"
	"time"

	"github.com/niakr1s/chatty-server/app/constants"
	"github.com/niakr1s/chatty-server/app/db"
	"github.com/niakr1s/chatty-server/app/events"
	log "github.com/sirupsen/logrus"
)

const (
	contentJSON = "application/json"
)

// Client ...
type Client struct {
	client *http.Client

	ctx context.Context

	Username string
	Password string
	URL      string
}

// New constructs a Bot
func New(ctx context.Context, username, password, url string) (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	return &Client{client: &http.Client{Jar: jar}, Username: username, Password: password, URL: url, ctx: ctx}, nil
}

// Connect trying to authorize at our chatty-server
func (b *Client) Connect() error {
	w, err := b.client.Post(b.URL+constants.RouteApi+constants.RouteAuthorize,
		contentJSON, strings.NewReader(fmt.Sprintf(`{"user": "%s", "password": "%s"}`, b.Username, b.Password)))
	if err != nil {
		return err
	}
	defer w.Body.Close()
	return nil
}

// StartListen polls events and sends them into channel.
// on network error closes that channel.
func (b *Client) StartListen() <-chan events.Event {
	ch := make(chan events.Event)
	go func() {
		pollURL := b.URL + constants.RouteApi + constants.RouteLoggedOnly + constants.RoutePoll
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

// PostMessage posts message in a new goroutine
func (b *Client) PostMessage(chat, text string) {
	log.Tracef("posting message in chat %s: %s", chat, text)
	go func() {
		postMsgURL := b.URL + constants.RouteApi + constants.RouteLoggedOnly + constants.RoutePostMessage
		if _, err := b.client.Post(postMsgURL, contentJSON,
			strings.NewReader(fmt.Sprintf(`{"user": "%s", "chat": "%s", "text": "%s"}`,
				b.Username, chat, text))); err != nil {
			log.Errorf("postMessage: %v", err)
		}
	}()
}

// GetChats ask server for all existent chats
func (b *Client) GetChats() ([]string, error) {
	getChatsURL := b.URL + constants.RouteApi + constants.RouteLoggedOnly + constants.RouteGetChats

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

// JoinChats joins chats, each join in a single goroutine
func (b *Client) JoinChats(chatnames ...string) {
	wg := sync.WaitGroup{}
	wg.Add(len(chatnames))
	for _, cn := range chatnames {
		go func(chatname string) {
			defer wg.Done()
			if err := b.JoinChat(chatname); err != nil {
				log.Error(err)
				return
			}
		}(cn)
	}
	wg.Wait()
}

// JoinChat joins a chat
func (b *Client) JoinChat(chat string) error {
	joinChatURL := b.URL + constants.RouteApi + constants.RouteLoggedOnly + constants.RouteJoinChat
	_, err := b.client.Post(joinChatURL, contentJSON, strings.NewReader(fmt.Sprintf(`{"chat": "%s"}`, chat)))
	if err != nil {
		return ErrJoinChat
	}
	return nil
}

// StartSendingKeepAlive sends keep-alive packages forever with interfal 10s
func (b *Client) StartSendingKeepAlive() {
	go func() {
		for {
			<-time.After(time.Second * 10)
			b.sendKeepAlive()
		}
	}()
}

// sendKeepAlive is a helper-function for startSendingKeepAlive
func (b *Client) sendKeepAlive() error {
	keepAliveURL := b.URL + constants.RouteApi + constants.RouteLoggedOnly + constants.RouteKeepAlive
	r, err := http.NewRequestWithContext(b.ctx, http.MethodPut, keepAliveURL, nil)
	if err != nil {
		return err
	}
	if _, err := b.client.Do(r); err != nil {
		return err
	}
	return nil
}
