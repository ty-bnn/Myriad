package parser

import (
	"myriad/tokenizer"
)

type Parser struct {
	tokens []tokenizer.Token
}

func (p *Parser)Parse(tokens []tokenizer.Token) (error) {
	p.tokens = tokens
	err := p.program(0)
	if err != nil {
		return err
	}

	return nil
}
