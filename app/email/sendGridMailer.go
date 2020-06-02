package email

import (
	"fmt"
	"net/smtp"
	"os"

	"github.com/mohamedattahri/mail"
	"github.com/niakr1s/chatty-server/app/constants"
	"github.com/niakr1s/chatty-server/app/er"
)

const sendGridURL = "smtp.sendgrid.net:25"

// SendGridMailer ...
type SendGridMailer struct {
	auth smtp.Auth
}

// NewSMTPMailer ...
func NewSMTPMailer() (*SendGridMailer, error) {
	res := &SendGridMailer{}
	apiKey := os.Getenv(constants.EnvSendGridAPIKey)
	if apiKey == "" {
		return nil, er.ErrEnvEmptySendGridAPIKey
	}
	res.auth = smtp.PlainAuth("", "apikey", apiKey, "smtp.sendgrid.net")

	if _, err := smtp.Dial(sendGridURL); err != nil {
		return nil, err
	}

	return res, nil
}

// SendMail ...
func (m *SendGridMailer) SendMail(email string, user string, activationToken string) error {
	msg, err := plainEmail(email, user, activationToken)
	if err != nil {
		return err
	}
	if err := smtp.SendMail(sendGridURL, m.auth, "chatsie", []string{email}, []byte(msg)); err != nil {
		return err
	}

	return nil
}

func plainEmail(email, user, activationToken string) (string, error) {
	from, err := mail.ParseAddress("pavel2188@gmail.com")
	if err != nil {
		return "", err
	}
	to, err := mail.ParseAddress(email)
	if err != nil {
		return "", err
	}
	msg := mail.NewMessage()
	msg.SetFrom(from)
	msg.To().Add(to)
	msg.SetSubject("Registration in Chatsie")
	msg.SetContentType("text/plain")
	fmt.Fprintf(msg.Body,
		`Hello, %[2]s!
You have succesfully registered at chatsie.herokuapp.com!
Here is your activation link: chatsie.herokuapp.com/api/verifyEmail/%[2]s/%[3]s`, email, user, activationToken)
	return msg.String(), nil
}
