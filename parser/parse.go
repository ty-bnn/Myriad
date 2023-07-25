package parser

import (
	"fmt"

	"myriad/tokenizer"
)

type Parser struct {
	tokens []tokenizer.Token
	index  int
}

func (p *Parser) Parse(tokens []tokenizer.Token) error {
	fmt.Println("Parseing...")
	p.tokens = tokens
	err := p.program()
	if err != nil {
		return err
	}

	fmt.Println("Parse Done.")
	return nil
}
