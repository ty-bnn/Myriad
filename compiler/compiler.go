package compiler

import (
	"myriad/tokenizer"
)

type Compiler struct {
	FunctionInterCodeMap *map[string][]InterCode
	FunctionVarMap *map[string][]Variable
	functionPointer string
	readFiles *[]string
}

func (c *Compiler) Compile(tokens []tokenizer.Token) error {
	c.FunctionInterCodeMap = &map[string][]InterCode{}
	c.FunctionVarMap = &map[string][]Variable{}
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
