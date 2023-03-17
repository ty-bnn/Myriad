package compiler

import (
	"dcc/tokenizer"
)

var functionInterCodeMap map[string][]InterCode
var functionVarMap map[string][]Variable
var functionPointer string
var readFiles []string

func Compile(tokens []tokenizer.Token) (map[string][]InterCode, map[string][]Variable, error) {
	functionInterCodeMap = map[string][]InterCode{}
	functionVarMap = map[string][]Variable{}
	functionPointer = "main"
	
	err := program(tokens, 0)
	if err != nil {
		return functionInterCodeMap, functionVarMap, err
	}

	// for debug.
	printInterCodes()

	return functionInterCodeMap, functionVarMap, err
}
