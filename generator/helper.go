package generator

import (
	"fmt"
	"errors"

	"dcc/types"
)

func generateCodeBlock(index int, functionName string, argValues []string) (int, []string, error) {
	var codes []string
	var err error
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
			if i, ok := getArgumentIndex(functionName, interCodes[index].Content); ok {
				code = argValues[i]
			} else if val, ok := argsInMain[interCodes[index].Content]; ok {
				code = val
			} else {
				return index, codes, errors.New(fmt.Sprintf("semantic error: variable is not declared"))
			}

			codes = append(codes, code)
			index++
		case types.CALLFUNC:
			_, codeBlock, err = generateCodeBlock(0, interCodes[index].Content, interCodes[index].ArgValues)
			if err != nil {
				return index, codes, err
			}

			codes = append(codes, codeBlock...)
			index++
		case types.IF:
			index, codeBlock, err = generateIfBlock(index, functionName, argValues)
			if err != nil {
				return index, codes, err
			}

			codes = append(codes, codeBlock...)
		case types.ENDIF:
			index++
			return index, codes, nil
		}
	}

	return index, codes, nil
}

func generateIfBlock(index int, functionName string, argValues []string) (int, []string, error) {
	var codes []string
	var err error

	interCodes := functionInterCodeMap[functionName]

	for index < len(interCodes) {
		if interCodes[index].Kind == types.IF || interCodes[index].Kind == types.ELIF {
			isTrue, err := getIfCondition(interCodes[index].IfContent, functionName, argValues)
			if err != nil {
				return index, codes, err
			}

			if isTrue {
				_, codes, err = generateCodeBlock(index + 1, functionName, argValues)
				if err != nil {
					return index, codes, err
				}

				return interCodes[index].IfContent.EndIndex, codes, nil
			}
		} else if interCodes[index].Kind == types.ELSE {
			index++
			index, codes, err = generateCodeBlock(index, functionName, argValues)
			if err != nil {
				return index, codes, err
			}

			return index, codes, nil
		}
		index++
	}

	return index, codes, nil
}
