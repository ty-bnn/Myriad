package compiler

import (
	"dcc/types"
)

var functionArgMap *map[string][]types.Argument
var functionCodeMap *map[string][]types.Code
var functionPointer string

func Compile(tokens []types.Token, faMap *map[string][]types.Argument, fcMap *map[string][]types.Code) (error) {
	functionArgMap = faMap
	functionCodeMap = fcMap
	
	err := program(tokens, 0)
	if err != nil {
		return err
	}

	return err
}
