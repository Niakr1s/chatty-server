package random

// StrGenerator ...
type StrGenerator interface {
	RandomStr() string
}

// StrGen ...
var StrGen StrGenerator

func init() {
	StrGen = NewRandStrGenerator()
}
