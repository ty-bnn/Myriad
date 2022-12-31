package generator

import (
	// "fmt"
	"dcc/types"
)

func generateCodeBlock(index int, functionName string, argValues []string) (int, []string) {
	var codes []string
	interCodes := functionInterCodeMap[functionName]
	
	for index < len(interCodes) {
		var code string
		var codeBlock []string

		switch interCodes[index].Kind {
		case types.ROW:
			code = interCodes[index].Content
			codes = append(codes, code)
			index++
		case types.VAR:
			i := getArgumentIndex(functionName, interCodes[index].Content)
			code = argValues[i]
			codes = append(codes, code)
			index++
		case types.CALLFUNC:
			_, codeBlock = generateCodeBlock(0, interCodes[index].Content, interCodes[index].ArgValues)
			codes = append(codes, codeBlock...)
			index++
		case types.IF:
			index, codeBlock = generateIfBlock(index, functionName, argValues)
			codes = append(codes, codeBlock...)
		case types.ENDIF:
			index++
			return index, codes
		}
	}

	return index, codes
}

func generateIfBlock(index int, functionName string, argValues []string) (int, []string) {
	var codes []string
	interCodes := functionInterCodeMap[functionName]

	for index < len(interCodes) {
		if interCodes[index].Kind == types.IF || interCodes[index].Kind == types.ELIF {
			if getIfCondition(interCodes[index].IfContent, functionName, argValues) {
				_, codes = generateCodeBlock(index + 1, functionName, argValues)
				return interCodes[index].IfContent.EndIndex, codes
			}
		} else if interCodes[index].Kind == types.ELSE {
			index++
			index, codes = generateCodeBlock(index, functionName, argValues)
			return index, codes
		}
		index++
	}

	return index, codes
}
