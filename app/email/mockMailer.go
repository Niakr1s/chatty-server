package email

import log "github.com/sirupsen/logrus"

// MockMailer used to mock Mailer interface
type MockMailer struct {
	Email          string
	User           string
	ActivationCode uint32
}

// NewMockMailer constructs MockMailer
func NewMockMailer() *MockMailer {
	return &MockMailer{}
}

// SendMail simply records input arguments in itself
func (m *MockMailer) SendMail(email string, user string, activationCode uint32) error {
	m.Email = email
	m.User = user
	m.ActivationCode = activationCode

	log.Debugf("activation link: /api/verifyEmail/%s/%d", user, activationCode)

	return nil
}
