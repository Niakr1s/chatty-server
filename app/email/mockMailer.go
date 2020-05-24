package email

import log "github.com/sirupsen/logrus"

// MockMailer used to mock Mailer interface
type MockMailer struct {
	Email           string
	User            string
	ActivationToken string
}

// NewMockMailer constructs MockMailer
func NewMockMailer() *MockMailer {
	return &MockMailer{}
}

// SendMail simply records input arguments in itself
func (m *MockMailer) SendMail(email string, user string, activationToken string) error {
	m.Email = email
	m.User = user
	m.ActivationToken = activationToken

	log.Debugf("activation link: /api/verifyEmail/%s/%s", user, activationToken)
	return nil
}
