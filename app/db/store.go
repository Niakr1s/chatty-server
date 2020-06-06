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
	models.User
	models.Chat
	Joined   bool                `json:"joined"`
	Messages []*models.Message   `json:"messages"`
	Users    []models.LoggedUser `json:"users"`
}

// NewChatReport ...
func NewChatReport(username, chatname string, joined bool) ChatReport {
	return ChatReport{
		User:     models.User{UserName: username},
		Chat:     models.Chat{ChatName: chatname},
		Joined:   joined,
		Messages: []*models.Message{},
		Users:    []models.LoggedUser{}}
}

// MakeChatReportForUser returns ChatReport for an user.
// If not joined, Messages and Users fields are always empty.
// Chat should be locked. LoggedDB, MessageDB should be unlocked.
func (s *Store) MakeChatReportForUser(username string, chat Chat) ChatReport {
	res := NewChatReport(username, chat.ChatName(), chat.IsInChat(username))

	if res.Joined {
		s.MessageDB.Lock()
		gotMessages, err := s.MessageDB.GetLastNMessages(chat.ChatName(), config.C.LastMessages)
		s.MessageDB.Unlock()
		if err == nil {
			res.Messages = gotMessages
		}
		s.LoggedDB.Lock()
		for _, u := range chat.GetUsers() {
			if loggedU, err := s.LoggedDB.Get(u.UserName); err == nil {
				res.Users = append(res.Users, *loggedU)
			}
		}
		s.LoggedDB.Unlock()
	}

	return res
}
