package bot

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"strings"

	"github.com/niakr1s/chatty-server/app/constants"
)

// Bot ...
type Bot struct {
	client *http.Client
}

// New ...
func New() (*Bot, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	return &Bot{client: &http.Client{Jar: jar}}, nil
}

// Connect ...
func (b *Bot) Connect(username, password, url string) error {
	w, err := b.client.Post(url+constants.RouteApi+constants.RouteAuthorize,
		"application/json", strings.NewReader(fmt.Sprintf(`{"user": "%s", "password": "%s"}`, username, password)))
	if err != nil {
		return err
	}
	w.Body.Close()
	return nil
}
