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
	Joined   bool                    `json:"joined"`
	Messages []*models.Message       `json:"messages"`
	Users    []models.UserWithStatus `json:"users"`
}

// NewChatReport ...
func NewChatReport(username, chatname string, joined bool) ChatReport {
	return ChatReport{
		User:     models.User{UserName: username},
		Chat:     models.Chat{ChatName: chatname},
		Joined:   joined,
		Messages: []*models.Message{},
		Users:    []models.UserWithStatus{}}
}

// MakeChatReportForUser returns ChatReport for an user.
// If not joined, Messages and Users fields are always empty.
// Chat should be locked. LoggedDB, MessageDB should be unlocked.
func (s *Store) MakeChatReportForUser(username string, chatname string) ChatReport {
	res := NewChatReport(username, chatname, s.ChatDB.IsInChat(chatname, username))

	if res.Joined {
		gotMessages, err := s.MessageDB.GetLastNMessages(chatname, config.C.LastMessages)
		if err == nil {
			res.Messages = gotMessages
		}
		for _, u := range s.ChatDB.GetUsers(chatname) {
			if loggedU, err := s.LoggedDB.Get(u); err == nil {
				userToAppend := models.NewUserWithStatus(loggedU.User, loggedU.UserStatus)
				res.Users = append(res.Users, userToAppend)
			}
		}
	}

	return res
}
