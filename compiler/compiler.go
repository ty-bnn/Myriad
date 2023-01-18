package compiler

import (
	"dcc/types"
)

var functionInterCodeMap map[string][]types.InterCode
var functionArgMap map[string][]types.Argument
var functionPointer string
var readFiles []string

func Compile(tokens []types.Token) (map[string][]types.InterCode, map[string][]types.Argument, error) {
	functionInterCodeMap = map[string][]types.InterCode{}
	functionArgMap = map[string][]types.Argument{}
	
	err := program(tokens, 0)
	if err != nil {
		return functionInterCodeMap, functionArgMap, err
	}

	return functionInterCodeMap, functionArgMap, err
}
