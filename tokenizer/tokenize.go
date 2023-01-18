package tokenizer

import (
	"dcc/types"
)

func Tokenize(lines []string) ([]types.Token, error) {
	var tokens []types.Token

	for i := 0; i < len(lines); i++ {
		var token types.Token
		var err error 
		index := 0
		for index < len(lines[i]) {
			switch lines[i][index] {
				case ' ':
					index++
					continue
				/*
				Read symbols '(', ')', ',', "[]", '{', '}', "==", "!="
				*/
				case '(', ')', ',', '[', '{', '}', '=', '!':
					index, token, err = readSymbols(index, lines[i], i)
				/*
				Read reserved words ("import", "from", "main", "if", "else if", "else")
				If not reserved words, read identifiers starts from 'i', 'f', 'm'.
				*/
				case 'e', 'f', 'i', 'm':
					index, token, err = readReservedWords(index, lines[i], i)
					if err != nil {
						index, token, err = readIdentifier(index, lines[i], i)
					}				
				/*
				Read Dfile commands.
				And Read Dfile arguments.
				*/
				case 'A', 'C', 'E', 'F', 'H', 'L', 'M', 'O', 'S', 'U', 'V', 'W':
					index, token, err = readDfCommands(index, lines[i], i)
					tokens = append(tokens, token)

					if err == nil {
						var dfArgTokens []types.Token
						index, dfArgTokens, err = readDfArgsPerLine(index, lines[i], i)
						tokens = append(tokens, dfArgTokens...)
						token = types.Token{Content: "\n", Kind: types.SDFARG}
					}
				case 'R':
					index, token, err = readDfCommands(index, lines[i], i)
					tokens = append(tokens, token)

					if err == nil {
						start := i
						for ;; {
							if i != start {
								tokens = append(tokens, types.Token{Content: "    ", Kind: types.SDFARG})
							}
							var dfArgTokens []types.Token
							index, dfArgTokens, err = readDfArgsPerLine(index, lines[i], i)
							tokens = append(tokens, dfArgTokens...)
							token = types.Token{Content: "\n", Kind: types.SDFARG}

							if err != nil || lines[i][len(lines[i])-1] != '\\' {
								break
							}

							tokens = append(tokens, token)
							i++
							index = 0
						}
					}
				/*
				Read strings start from " and ends at ".
				*/
				case '"':
					index, token, err = readString(index, lines[i], i)
				/*
				Read identifiers.
				*/
				default:
					index, token, err = readIdentifier(index, lines[i], i)
			}

			if err != nil {
				return nil, err
			}

			tokens = append(tokens, token)
		}
	}

	return tokens, nil
}