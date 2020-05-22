package random

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandGenerator_RandomString(t *testing.T) {
	r := NewRandStrGenerator()

	s1 := r.RandomStr()
	s2 := r.RandomStr()

	assert.NotEqual(t, s1, s2)
}
