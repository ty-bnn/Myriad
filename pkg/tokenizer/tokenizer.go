package tokenizer

import "github.com/ty-bnn/myriad/pkg/model/token"

type Tokenizer struct {
	lines  []string
	Tokens []token.Token
}

func New(lines []string) *Tokenizer {
	return &Tokenizer{lines: lines}
}
