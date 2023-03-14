package tokenizer

import (
	"fmt"
)

func Tokenize(lines []string) ([]Token, error) {
	fmt.Println("Tokenize...")

	var tokens []Token
	var err error
	var command string

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
						if err == nil {
							command = newToken.Content
							isInCommand = true
						} else {
							i, newToken, err = readIdentifier(i, line, row)
						}

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
				var space string

				i, newTokens, err = readDfArgs(i, line, row)

				if (0 <= len(tokens) - 1 && tokens[len(tokens) - 1].Kind == SDFCOMMAND) {
					// 既登録のトークン列の末尾がDfコマンドの場合，スペースを一つ空ける
					space = " "	
				} else {
					// 既登録のトークン列の末尾がDfコマンドではない場合，スペースをDfコマンド長+1分空ける
					for j := 0; j <= len(command); j++ {
						space = space + " "
					}
				}
				tokens = append(tokens, Token{Content: space, Kind: SDFARG})

				tokens = append(tokens, newTokens...)
				tokens = append(tokens, Token{Content: "\n", Kind: SDFARG})
				

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
