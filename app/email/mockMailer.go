package email

import log "github.com/sirupsen/logrus"

// MockMailer used to mock Mailer interface
type MockMailer struct {
	Email              string
	User               string
	ActivationToken    string
	ResetPasswordToken string
}

// NewMockMailer constructs MockMailer
func NewMockMailer() *MockMailer {
	return &MockMailer{}
}

// SendActivationEmail simply records input arguments in itself
func (m *MockMailer) SendActivationEmail(email string, user string, activationToken string) error {
	m.Email = email
	m.User = user
	m.ActivationToken = activationToken

	log.Infof("activation link: /api/verifyEmail/%s/%s", user, activationToken)
	return nil
}

// SendResetPasswordEmail simply records input arguments in itself
func (m *MockMailer) SendResetPasswordEmail(email string, user string, resetPasswordToken string) error {
	m.Email = email
	m.User = user
	m.ResetPasswordToken = resetPasswordToken

	log.Infof("reset password token: %s", resetPasswordToken)
	return nil
}
