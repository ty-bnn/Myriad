package generator

import (
	"dcc/compiler"
)

var mainCodes []compiler.InterCode
var argsInMain map[string]string
var command string

func GenerateCode(fInterCodeMap map[string][]compiler.InterCode, fArgMap map[string][]compiler.Variable) ([]string, error) {
	argsInMain = map[string]string{}
	mainCodes = fInterCodeMap["main"]
	command = ""

	_, codes, err := generateCodeBlock(0)
	if err != nil {
		return codes, err
	}

	return codes, nil
}
