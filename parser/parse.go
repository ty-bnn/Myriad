package parser

import (
	"myriad/tokenizer"
)

type Parser struct {
	tokens []tokenizer.Token
	index int
}

func (p *Parser)Parse(tokens []tokenizer.Token) error {
	p.tokens = tokens
	err := p.program()
	if err != nil {
		return err
	}

	return nil
}
