package generator

import (
	"fmt"

	"dcc/compiler"
)

func generateCodeBlock(index int, functionName string) (int, []string, error) {
	fmt.Println("---------- generating code ---------")

	var codes []string
	var err error
	interCodes := functionInterCodeMap[functionName]
	
	for index < len(interCodes) {
		var code string
		var codeBlock []string

		switch interCodes[index].Kind {
		case compiler.ROW:
			code = interCodes[index].Content
			codes = append(codes, code)
			index++
		case compiler.COMMAND:
			if command == "RUN" && interCodes[index].Content == "RUN" {
				// RUN命令の結合
				codes[len(codes) - 1] = codes[len(codes) - 1][:len(codes[len(codes) - 1]) - 1] + " \\\n"
				code = "    "
				index++
			} else {
				code = interCodes[index].Content
				command = interCodes[index].Content
			}
			codes = append(codes, code)
			index++
		case compiler.IF:
			index, codeBlock, err = generateIfBlock(index, functionName)
			if err != nil {
				return index, codes, err
			}

			codes = append(codes, codeBlock...)
		case compiler.ENDIF:
			index++
			return index, codes, nil
		}
	}

	return index, codes, nil
}

func generateIfBlock(index int, functionName string) (int, []string, error) {
	var codes []string
	var err error

	interCodes := functionInterCodeMap[functionName]

	for index < len(interCodes) {
		if interCodes[index].Kind == compiler.IF || interCodes[index].Kind == compiler.ELIF {
			isTrue, err := getIfCondition(interCodes[index].IfContent, functionName)
			if err != nil {
				return index, codes, err
			}
			
			if isTrue {
				_, codes, err = generateCodeBlock(index + 1, functionName)
				if err != nil {
					return index, codes, err
				}

				return index + interCodes[index].IfContent.EndIndex, codes, nil
			}
		} else if interCodes[index].Kind == compiler.ELSE {
			index++
			index, codes, err = generateCodeBlock(index, functionName)
			if err != nil {
				return index, codes, err
			}

			return index, codes, nil
		} else if interCodes[index].Kind == compiler.ENDIF {
			index++
			break
		}

		index++
	}

	return index, codes, nil
}
