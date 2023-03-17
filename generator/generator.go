package generator

import (
	"dcc/compiler"
)

var functionInterCodeMap map[string][]compiler.InterCode
var functionVarMap map[string][]compiler.Variable
var argsInMain map[string]string
var command string

func GenerateCode(fInterCodeMap map[string][]compiler.InterCode, fArgMap map[string][]compiler.Variable) ([]string, error) {
	functionInterCodeMap = fInterCodeMap
	functionVarMap = fArgMap
	argsInMain = map[string]string{}
	command = ""

	_, codes, err := generateCodeBlock(0, "main")
	if err != nil {
		return codes, err
	}

	return codes, nil
}
