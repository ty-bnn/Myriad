package generator

import (
	"github.com/ty-bnn/myriad/pkg/model/codes"
)

type Generator struct {
	funcToCodes map[string][]codes.Code
	RawCodes    []string
	funcPtr     string
	index       int
	commandPtr  string
}

func NewGenerator(funcToCodes map[string][]codes.Code) *Generator {
	return &Generator{
		funcToCodes: funcToCodes,
	}
}
