package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"

	"github.com/niakr1s/chatty-server/app/events"
)

var (
	errUnknownEvent = errors.New("unknown event")
)

func parseEvent(r io.ReadCloser) (events.Event, error) {
	b, err := ioutil.ReadAll(r)
	r.Close()
	if err != nil {
		return nil, err
	}

	type typeExtractor struct {
		Type string `json:"type"`
	}
	t := typeExtractor{}
	if err := json.NewDecoder(bytes.NewReader(b)).Decode(&t); err != nil {
		return nil, err
	}

	ewt := events.EventWithType{}
	switch t.Type {
	case "ChatJoinEvent":
		ewt.Event = &events.ChatJoinEvent{}
	case "ChatLeaveEvent":
		ewt.Event = &events.ChatLeaveEvent{}
	case "LoginEvent":
		ewt.Event = &events.LoginEvent{}
	case "LogoutEvent":
		ewt.Event = &events.LogoutEvent{}
	case "ChatCreatedEvent":
		ewt.Event = &events.ChatCreatedEvent{}
	case "ChatRemovedEvent":
		ewt.Event = &events.ChatRemovedEvent{}
	case "MessageEvent":
		ewt.Event = &events.MessageEvent{}
	default:
		return nil, errUnknownEvent
	}

	if err := json.NewDecoder(bytes.NewReader(b)).Decode(&ewt); err != nil {
		return nil, err
	}
	return ewt.Event, nil
}
