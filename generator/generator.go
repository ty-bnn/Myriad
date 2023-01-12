package generator

import (
	"os"

	"dcc/types"
)

var functionInterCodeMap map[string][]types.InterCode
var functionArgMap map[string][]types.Argument

func GenerateCode(fInterCodeMap map[string][]types.InterCode, fArgMap map[string][]types.Argument) []string {
	functionInterCodeMap = fInterCodeMap
	functionArgMap = fArgMap

	_, codes := generateCodeBlock(0, "main", os.Args[3:])

	return codes
}
