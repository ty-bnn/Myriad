package generator

import (
	"fmt"
	"os"
	"errors"

	"dcc/types"
)

var functionInterCodeMap map[string][]types.InterCode
var functionArgMap map[string][]types.Argument
var argsInMain map[string]string

func GenerateCode(fInterCodeMap map[string][]types.InterCode, fArgMap map[string][]types.Argument) ([]string, error) {
	functionInterCodeMap = fInterCodeMap
	functionArgMap = fArgMap
	argsInMain = map[string]string{}

	if len(functionArgMap["main"]) != len(os.Args[3:]) {
		return []string{}, errors.New(fmt.Sprintf("system error: length of main argument is not match"))
	}

	for i := 0; i < len(functionArgMap["main"]); i++ {
		argsInMain[functionArgMap["main"][i].Name] = os.Args[i + 3]
	}

	_, codes, err := generateCodeBlock(0, "main", os.Args[3:])
	if err != nil {
		return codes, err
	}

	return codes, nil
}
