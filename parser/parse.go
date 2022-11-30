package parser

import (
	"dcc/tokenizer"
)

func Parse(tokens []tokenizer.Token) (error) {
	err := program(tokens, 0)

	return err
}