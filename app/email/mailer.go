package email

// Mailer is used to send email to user with activationCode
type Mailer interface {
	SendMail(email string, user string, activationToken string) error
}
