package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmail_GenerateActivationToken(t *testing.T) {
	email := Email{Address: "user@example.org"}

	email.GenerateActivationToken()
	token1 := email.ActivationToken

	email.GenerateActivationToken()
	token2 := email.ActivationToken

	assert.NotEqual(t, token1, token2)
}
