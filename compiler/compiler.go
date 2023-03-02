package compiler

import (
	"dcc/tokenizer"
)

var functionInterCodeMap map[string][]InterCode
var functionArgMap map[string][]Argument
var functionPointer string
var readFiles []string

func Compile(tokens []tokenizer.Token) (map[string][]InterCode, map[string][]Argument, error) {
	functionInterCodeMap = map[string][]InterCode{}
	functionArgMap = map[string][]Argument{}
	
	err := program(tokens, 0)
	if err != nil {
		return functionInterCodeMap, functionArgMap, err
	}

	return functionInterCodeMap, functionArgMap, err
}
