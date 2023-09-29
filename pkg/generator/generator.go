package generator

import (
	"github.com/ty-bnn/myriad/pkg/model/codes"
)

type Generator struct {
	FuncToCodes map[string][]codes.Code
	Codes       []string
	command     string
}

func New(funcToCodes map[string][]codes.Code) *Generator {
	return &Generator{
		FuncToCodes: funcToCodes,
	}
}
