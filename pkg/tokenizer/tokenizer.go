package tokenizer

import "github.com/ty-bnn/myriad/pkg/model/token"

type Tokenizer struct {
	data        string
	p           int
	isInDfBlock bool
	isInCommand bool
	commandPtr  string
	Tokens      []token.Token
}

func New(data string) *Tokenizer {
	return &Tokenizer{data: data}
}
