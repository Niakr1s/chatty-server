package db

import (
	"github.com/niakr1s/chatty-server/app/config"
	"github.com/niakr1s/chatty-server/app/models"
)

// Store contains all databases
type Store struct {
	UserDB    UserDB
	ChatDB    ChatDB
	LoggedDB  LoggedDB
	MessageDB MessageDB
}

// NewStore ...
func NewStore(u UserDB, c ChatDB, l LoggedDB, m MessageDB) *Store {
	return &Store{UserDB: u, ChatDB: c, LoggedDB: l, MessageDB: m}
}

// ChatReport ...
type ChatReport struct {
	models.Chat
	Joined   bool              `json:"joined"`
	Messages []*models.Message `json:"messages"`
	Users    []models.User     `json:"users"`
}

// NewChatReport ...
func NewChatReport(chatname string, joined bool) ChatReport {
	return ChatReport{Chat: models.Chat{ChatName: chatname}, Joined: joined, Messages: []*models.Message{}, Users: []models.User{}}
}

// MakeChatReportForUser returns ChatReport for an user
// if not joined, Messages and Users fields are always empty
func (s *Store) MakeChatReportForUser(username string, chatname string) (ChatReport, error) {
	c, err := s.ChatDB.Get(chatname)
	if err != nil {
		return ChatReport{}, err
	}

	c.Lock()
	defer c.Unlock()

	res := NewChatReport(chatname, c.IsInChat(username))

	if res.Joined {
		gotMessages, err := s.MessageDB.GetLastNMessages(c.ChatName(), config.C.LastMessages)
		if err == nil {
			res.Messages = gotMessages
		}
		res.Users = c.GetUsers()
	}

	return res, nil
}
