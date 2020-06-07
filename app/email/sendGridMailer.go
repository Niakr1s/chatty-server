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

// SendActivationEmail ...
func (m *SendGridMailer) SendActivationEmail(email string, user string, activationToken string) error {
	msg, err := activationBody(email, user, activationToken)
	if err != nil {
		return err
	}
	if err := smtp.SendMail(sendGridURL, m.auth, "chatsie", []string{email}, []byte(msg)); err != nil {
		return err
	}

	return nil
}

func activationBody(email, user, activationToken string) (string, error) {
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

// SendResetPasswordEmail ...
func (m *SendGridMailer) SendResetPasswordEmail(email string, user string, resetPasswordToken string) error {
	msg, err := resetPasswordBody(email, user, resetPasswordToken)
	if err != nil {
		return err
	}
	if err := smtp.SendMail(sendGridURL, m.auth, "chatsie", []string{email}, []byte(msg)); err != nil {
		return err
	}

	return nil
}

func resetPasswordBody(email, user, resetPasswordToken string) (string, error) {
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
	msg.SetSubject("Password reset in Chatsie")
	msg.SetContentType("text/plain")
	fmt.Fprintf(msg.Body,
		`Hello, %s!
You have requested password reset at chatsie.herokuapp.com!
Here is your password reset token, you must submit it at password reset form: %s`, user, resetPasswordToken)
	return msg.String(), nil
}
