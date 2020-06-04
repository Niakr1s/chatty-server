package email

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSendGridMailer(t *testing.T) {
	os.Setenv("SENDGRID_KEY", "123456")

	_, err := NewSMTPMailer()

	assert.NoError(t, err)
}
