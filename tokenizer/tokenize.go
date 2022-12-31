package tokenizer

import (
	"dcc/types"
)

func Tokenize(lines []string) ([]types.Token, error) {
	var tokens []types.Token

	for line, lineStr := range lines {
		var token types.Token
		var err error 
		index := 0
		for index < len(lineStr) {
			switch lineStr[index] {
				case ' ':
					index++
					continue
				/*
				Read symbols '(', ')', ',', "[]", '{', '}', "==", "!="
				*/
				case '(', ')', ',', '[', '{', '}', '=', '!':
					index, token, err = readSymbols(index, lineStr, line)
				/*
				Read reserved words ("import", "from", "main", "if", "else if", "else")
				If not reserved words, read identifiers starts from 'i', 'f', 'm'.
				*/
				case 'e', 'f', 'i', 'm':
					index, token, err = readReservedWords(index, lineStr, line)
					if err != nil {
						index, token, err = readIdentifier(index, lineStr, line)
					}				
				/*
				Read Dfile commands.
				And Read Dfile arguments.
				*/
				case 'A', 'C', 'E', 'F', 'H', 'L', 'M', 'O', 'R', 'S', 'U', 'V', 'W':
					index, token, err = readDfCommands(index, lineStr, line)
					tokens = append(tokens, token)

					if err == nil {
						for index < len(lineStr) {
							// for index < len(lineStr) {
							// 	if lineStr[index] != ' ' {
							// 		break
							// 	}
							// 	index++
							// }
							index, token, err = readDfArgs(index, lineStr, line)
							if err != nil {
								break
							}

							tokens = append(tokens, token)
						}
						token = types.Token{Content: "\n", Kind: types.SDFARG}
					}
				/*
				Read strings start from " and ends at ".
				*/
				case '"':
					index, token, err = readString(index, lineStr, line)
				/*
				Read identifiers.
				*/
				default:
					index, token, err = readIdentifier(index, lineStr, line)
			}

			if err != nil {
				return nil, err
			}

			tokens = append(tokens, token)
		}
	}

	return tokens, nil
}