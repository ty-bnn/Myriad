package tokenizer

import (
	"fmt"
)

func Tokenize(lines []string) ([]Token, error) {
	fmt.Println("Tokenize...")

	var tokens []Token
	var err error

	isInCommand := false

	for row, line := range lines {
		i := 0

		for i < len(line) {
			if line[i] == ' ' {
				i++
				continue
			}

			if !isInCommand {
				// Dockerfileの引数以外
				var newToken Token
				switch line[i] {
					/*
					Read Dfile commands.
					And Read Dfile arguments.
					*/
					case 'A', 'C', 'E', 'F', 'H', 'L', 'M', 'O', 'R', 'S', 'U', 'V', 'W':
						i, newToken, err = readDfCommands(i, line, row)
						if err != nil {
							i, newToken, err = readIdentifier(i, line, row)
						}
						isInCommand = true

					/*
					Read reserved words ("import", "from", "main", "if", "else if", "else")
					If not reserved words, read identifiers starts from 'i', 'f', 'm'.
					*/
					case 'e', 'f', 'i', 'm':
						i, newToken, err = readReservedWords(i, line, row)
						if err != nil {
							i, newToken, err = readIdentifier(i, line, row)
						}

					/*
					Read strings start from " and ends at ".
					*/
					case '"':
						i, newToken, err = readString(i, line, row)

					/*
					Read symbols '(', ')', ',', "[]", '{', '}', "==", "!=".
					*/
					case '(', ')', ',', '[', '{', '}', '=', '!':
						i, newToken, err = readSymbols(i, line, row)
					
					/*
					Read identifier.
					*/
					default:
						i, newToken, err = readIdentifier(i, line, row)
				}
				tokens = append(tokens, newToken)
			} else {
				// Dockerfileの引数
				var newTokens []Token
				i, newTokens, err = readDfArgs(i, line, row)
				tokens = append(tokens, newTokens...)

				if line[len(line) - 1] != '\\' {
					isInCommand = false
				}
			}

			if err != nil {
				return tokens, err
			}
		}
	}

	// for Debug.
	printTokens(tokens)

	return tokens, nil
}
