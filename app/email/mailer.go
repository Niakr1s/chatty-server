package email

// Mailer is used to send email to user with activationCode
type Mailer interface {
	SendActivationEmail(email string, user string, activationToken string) error
	SendResetPasswordEmail(email string, user string, resetPasswordToken string) error
}
