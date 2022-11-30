package compiler

import (
	"dcc/tokenizer"
)

var dfCodes []string

func Generate(tokens []tokenizer.Token) ([]string, error) {
	err := program(tokens, 0)

	return dfCodes, err
}
