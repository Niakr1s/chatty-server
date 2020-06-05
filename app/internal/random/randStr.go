package random

import (
	"fmt"
	"math/rand"
	"time"
)

// StrGenerator ...
type StrGenerator interface {
	RandomStr() string
}

// StrGen ...
var StrGen StrGenerator

func init() {
	StrGen = NewRandStrGenerator()
}

// RandStrGenerator ...
type RandStrGenerator struct {
	r *rand.Rand
}

// NewRandStrGenerator ...
func NewRandStrGenerator() *RandStrGenerator {
	res := &RandStrGenerator{}
	res.r = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	return res
}

// RandomStr ...
func (r *RandStrGenerator) RandomStr() string {
	u := r.r.Uint64()

	str := fmt.Sprintf("%d", u)
	return str
}
