package compiler

import (
	"fmt"
	"myriad/tokenizer"
)

type Compiler struct {
	tokens               []tokenizer.Token
	FunctionInterCodeMap map[string][]InterCode
	FunctionVarMap       map[string][]Variable
	functionPointer      string
	readFiles            []string
	index                int
}

func (c *Compiler) Compile(tokens []tokenizer.Token) error {
	fmt.Println("Compiling...")
	c.FunctionInterCodeMap = make(map[string][]InterCode)
	c.FunctionVarMap = make(map[string][]Variable)
	c.functionPointer = "main"

	err := c.program(tokens, 0)
	if err != nil {
		return err
	}

	fmt.Println("Compile Done.")

	return err
}
