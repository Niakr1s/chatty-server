package models

import (
	"math/rand"
)

// Email ...
type Email struct {
	Address         string `json:"email" validate:"email"`
	ActivationToken string `json:"-" validate:"required"`
	Activated       bool
}

// GenerateActivationToken ...
func (e *Email) GenerateActivationToken() {
	e.ActivationToken = randSeq(20)
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
