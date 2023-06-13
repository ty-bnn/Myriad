package compiler

import (
	"myriad/tokenizer"
)

type Compiler struct {
	tokens []tokenizer.Token
	FunctionInterCodeMap map[string][]InterCode
	FunctionVarMap map[string][]Variable
	functionPointer string
	readFiles []string
	index int
}

func (c *Compiler) Compile(tokens []tokenizer.Token) error {
	c.FunctionInterCodeMap = make(map[string][]InterCode)
	c.FunctionVarMap = make(map[string][]Variable)
	c.functionPointer = "main"
	
	err := c.program(tokens, 0)
	if err != nil {
		return err
	}

	// for debug.
	c.printInterCodes("main")
	c.printInterCodes("abc")

	return err
}
