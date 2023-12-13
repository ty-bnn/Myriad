package parser

import (
	"github.com/ty-bnn/myriad/pkg/model/token"

	"github.com/ty-bnn/myriad/pkg/model/codes"
)

type Parser struct {
	tokens        []token.Token
	FuncToCodes   map[string][]codes.Code
	compiledFiles []string
	index         int
	filePath      string
}

func NewParser(tokens []token.Token, filePath string) *Parser {
	return &Parser{
		tokens:      tokens,
		FuncToCodes: make(map[string][]codes.Code),
		filePath:    filePath,
	}
}
