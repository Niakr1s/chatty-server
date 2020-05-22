package random

import (
	"fmt"
	"math/rand"
	"time"
)

// RandStrGenerator ...
type RandStrGenerator struct {
	r *rand.Rand
}

// NewRandStrGenerator ...
func NewRandStrGenerator() *RandStrGenerator {
	res := &RandStrGenerator{}
	res.r = rand.New(rand.NewSource(time.Now().UnixNano()))
	return res
}

// RandomStr ...
func (r *RandStrGenerator) RandomStr() string {
	u := r.r.Uint64()

	str := fmt.Sprintf("%d", u)
	return str
}
