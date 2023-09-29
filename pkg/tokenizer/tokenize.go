package tokenizer

import (
	"fmt"

	"github.com/ty-bnn/myriad/pkg/model/token"
)

func (t *Tokenizer) Tokenize() error {
	fmt.Println("Tokenizing...")

	var err error
	var command string

	isInCommand := false

	for row, line := range t.lines {
		i := 0

		for i < len(line) {
			if line[i] == ' ' {
				i++
				continue
			}

			if !isInCommand {
				// Dockerfileの引数以外
				var newToken token.Token
				switch line[i] {
				case 'A', 'C', 'E', 'F', 'H', 'L', 'M', 'O', 'R', 'S', 'U', 'V', 'W':
					i, newToken, err = readDfCommands(i, line, row)
					if err == nil {
						command = newToken.Content
						isInCommand = true
					} else {
						i, newToken, err = readIdentifier(i, line, row)
					}
				case 'e', 'f', 'i', 'm':
					i, newToken, err = readReservedWords(i, line, row)
					if err != nil {
						i, newToken, err = readIdentifier(i, line, row)
					}
				case '"':
					i, newToken, err = readString(i, line, row)
				case '(', ')', ',', '[', ']', '{', '}', '=', '!', ':':
					i, newToken, err = readSymbols(i, line, row)
				case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
					i, newToken, err = readNumber(i, line, row)
				default:
					i, newToken, err = readIdentifier(i, line, row)
				}
				t.Tokens = append(t.Tokens, newToken)
			} else {
				// Dockerfileの引数
				var newTokens []token.Token
				var space string

				i, newTokens, err = readDfArgs(i, line, row)

				if 0 <= len(t.Tokens)-1 && t.Tokens[len(t.Tokens)-1].Kind == token.DFCOMMAND {
					// 既登録のトークン列の末尾がDfコマンドの場合，スペースを一つ空ける
					space = " "
				} else {
					// 既登録のトークン列の末尾がDfコマンドではない場合，スペースをDfコマンド長+1分空ける
					for j := 0; j <= len(command); j++ {
						space = space + " "
					}
				}
				t.Tokens = append(t.Tokens, token.Token{Content: space, Kind: token.DFARG})

				t.Tokens = append(t.Tokens, newTokens...)
				t.Tokens = append(t.Tokens, token.Token{Content: "\n", Kind: token.DFARG})

				if line[len(line)-1] != '\\' {
					isInCommand = false
				}
			}

			if err != nil {
				return err
			}
		}
	}

	fmt.Println("Tokenize Done.")

	return nil
}
