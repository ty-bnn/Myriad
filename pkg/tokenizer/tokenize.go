package tokenizer

import (
	"errors"
	"fmt"
	. "strings"

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
				case 'e', 'f', 'i', 'k', 'm', 'v', 'J':
					i, newToken, err = readReservedWords(i, line, row)
					if err != nil {
						i, newToken, err = readIdentifier(i, line, row)
					}
				case '"':
					i, newToken, err = readString(i, line, row)
				case '(', ')', ',', '[', ']', '{', '}', '=', '!', ':', '.', '&', '|':
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

// isAlphabetOrNumber judges whether the char is alphabet or number
// 予約語の後ろにアルファベットや数字が続く場合は予約語ではなく識別子であると判断する
func isAlphabetOrNumber(line string, index int) bool {
	if len(line) <= index || (line[index] < '0' || ('9' < line[index] && line[index] < 'A') || ('Z' < line[index] && line[index] < 'a') || ('z' < line[index])) {
		return false
	}

	return true
}

func readReservedWords(index int, line string, row int) (int, token.Token, error) {
	if index+6 <= len(line) && line[index:index+6] == "import" && !isAlphabetOrNumber(line, index+6) {
		return index + 6, token.Token{Content: "import", Kind: token.IMPORT, Line: row + 1}, nil
	} else if index+4 <= len(line) && line[index:index+4] == "from" && !isAlphabetOrNumber(line, index+4) {
		return index + 4, token.Token{Content: "from", Kind: token.FROM, Line: row + 1}, nil
	} else if index+4 <= len(line) && line[index:index+4] == "main" && !isAlphabetOrNumber(line, index+4) {
		return index + 4, token.Token{Content: "main", Kind: token.MAIN, Line: row + 1}, nil
	} else if index+2 <= len(line) && line[index:index+2] == "if" && !isAlphabetOrNumber(line, index+2) {
		return index + 2, token.Token{Content: "if", Kind: token.IF, Line: row + 1}, nil
	} else if index+7 <= len(line) && line[index:index+7] == "else if" && !isAlphabetOrNumber(line, index+7) {
		return index + 7, token.Token{Content: "else if", Kind: token.ELIF, Line: row + 1}, nil
	} else if index+4 <= len(line) && line[index:index+4] == "else" && !isAlphabetOrNumber(line, index+4) {
		return index + 4, token.Token{Content: "else", Kind: token.ELSE, Line: row + 1}, nil
	} else if index+3 <= len(line) && line[index:index+3] == "for" && !isAlphabetOrNumber(line, index+3) {
		return index + 3, token.Token{Content: "for", Kind: token.FOR, Line: row + 1}, nil
	} else if index+2 <= len(line) && line[index:index+2] == "in" && !isAlphabetOrNumber(line, index+2) {
		return index + 2, token.Token{Content: "in", Kind: token.IN, Line: row + 1}, nil
	} else if index+4 <= len(line) && line[index:index+4] == "keys" && !isAlphabetOrNumber(line, index+4) {
		return index + 4, token.Token{Content: "keys", Kind: token.KEYS, Line: row + 1}, nil
	} else if index+13 <= len(line) && line[index:index+13] == "JsonUnmarshal" && !isAlphabetOrNumber(line, index+13) {
		return index + 13, token.Token{Content: "JsonUnmarshal", Kind: token.JSONUNMARSHAL, Line: row + 1}, nil
	}

	return index, token.Token{}, errors.New(fmt.Sprintf("tokenize error: found invalid character in line %d", row))
}

func readSymbols(index int, line string, row int) (int, token.Token, error) {
	if line[index] == '(' {
		return index + 1, token.Token{Content: "(", Kind: token.LPAREN, Line: row + 1}, nil
	} else if line[index] == ')' {
		return index + 1, token.Token{Content: ")", Kind: token.RPAREN, Line: row + 1}, nil
	} else if line[index] == ',' {
		return index + 1, token.Token{Content: ",", Kind: token.COMMA, Line: row + 1}, nil
	} else if line[index] == '[' {
		return index + 1, token.Token{Content: "[", Kind: token.LBRACKET, Line: row + 1}, nil
	} else if line[index] == ']' {
		return index + 1, token.Token{Content: "]", Kind: token.RBRACKET, Line: row + 1}, nil
	} else if line[index] == '{' {
		return index + 1, token.Token{Content: "{", Kind: token.LBRACE, Line: row + 1}, nil
	} else if line[index] == '}' {
		return index + 1, token.Token{Content: "}", Kind: token.RBRACE, Line: row + 1}, nil
	} else if line[index] == '.' {
		return index + 1, token.Token{Content: ".", Kind: token.DOT, Line: row + 1}, nil
	} else if index+2 <= len(line) && line[index:index+2] == "==" {
		return index + 2, token.Token{Content: "==", Kind: token.EQUAL, Line: row + 1}, nil
	} else if index+2 <= len(line) && line[index:index+2] == "!=" {
		return index + 2, token.Token{Content: "!=", Kind: token.NOTEQUAL, Line: row + 1}, nil
	} else if index+2 <= len(line) && line[index:index+2] == ":=" {
		return index + 2, token.Token{Content: ":=", Kind: token.DEFINE, Line: row + 1}, nil
	} else if line[index] == '=' {
		return index + 1, token.Token{Content: "=", Kind: token.ASSIGN, Line: row + 1}, nil
	} else if index+2 <= len(line) && line[index:index+2] == "&&" {
		return index + 2, token.Token{Content: "&&", Kind: token.AND, Line: row + 1}, nil
	} else if index+2 <= len(line) && line[index:index+2] == "||" {
		return index + 2, token.Token{Content: "||", Kind: token.OR, Line: row + 1}, nil
	}

	return index, token.Token{}, errors.New(fmt.Sprintf("Symbols'index: %d find invalid character in \"%s\".", index, line))
}

func readNumber(index int, line string, row int) (int, token.Token, error) {
	begin := index
	for index < len(line) && '0' <= line[index] && line[index] <= '9' {
		index++
	}

	return index, token.Token{Content: line[begin:index], Kind: token.NUMBER, Line: row + 1}, nil

}

func readDfCommands(index int, line string, row int) (int, token.Token, error) {
	/*
		ADD, ARG, CMD, COPY, ENTRYPOINT, ENV, EXPOSE, FROM,
		HEALTHCHECK, LABEL, MAINTAINER, ONBUILD, RUN, SHELL,
		STOPSIGNAL, USER, VOLUME, WORKDIR
	*/
	if HasPrefix(line[index:], "ADD") || HasPrefix(line[index:], "ARG") || HasPrefix(line[index:], "CMD") || HasPrefix(line[index:], "ENV") || HasPrefix(line[index:], "RUN") {
		return index + 3, token.Token{Content: line[index : index+3], Kind: token.DFCOMMAND, Line: row + 1}, nil
	}

	if HasPrefix(line[index:], "COPY") || HasPrefix(line[index:], "FROM") || HasPrefix(line[index:], "USER") {
		return index + 4, token.Token{Content: line[index : index+4], Kind: token.DFCOMMAND, Line: row + 1}, nil
	}

	if HasPrefix(line[index:], "LABEL") || HasPrefix(line[index:], "SHELL") {
		return index + 5, token.Token{Content: line[index : index+5], Kind: token.DFCOMMAND, Line: row + 1}, nil
	}

	if HasPrefix(line[index:], "") || HasPrefix(line[index:], "") {
		return index + 6, token.Token{Content: line[index : index+6], Kind: token.DFCOMMAND, Line: row + 1}, nil
	}

	if HasPrefix(line[index:], "ONBUILD") || HasPrefix(line[index:], "WORKDIR") {
		return index + 7, token.Token{Content: line[index : index+7], Kind: token.DFCOMMAND, Line: row + 1}, nil
	}

	if HasPrefix(line[index:], "ENTRYPOINT") || HasPrefix(line[index:], "MAINTAINER") || HasPrefix(line[index:], "STOPSIGNAL") {
		return index + 10, token.Token{Content: line[index : index+10], Kind: token.DFCOMMAND, Line: row + 1}, nil
	}

	if HasPrefix(line[index:], "HEALTHCHECK") {
		return index + 11, token.Token{Content: line[index : index+10], Kind: token.DFCOMMAND, Line: row + 1}, nil
	}

	return index, token.Token{}, errors.New(fmt.Sprintf("tokenize error: found invalid character in line %d", row))
}

func readDfArg(index int, line string, row int) (int, []token.Token, error) {
	var tokens []token.Token

	if index+2 < len(line) && line[index:index+2] == "{{" {
		// {{
		tokens = append(tokens, token.Token{Content: "{{", Kind: token.LDOUBLEBRA, Line: row + 1})
		index += 2
		for index < len(line)-1 {
			if line[index:index+2] == "}}" {
				tokens = append(tokens, token.Token{Content: "}}", Kind: token.RDOUBLEBRA, Line: row + 1})
				index = index + 2
				return index, tokens, nil
			} else if line[index] == '[' {
				tokens = append(tokens, token.Token{Content: "[", Kind: token.LBRACKET, Line: row + 1})
				index++
			} else if line[index] == ']' {
				tokens = append(tokens, token.Token{Content: "]", Kind: token.RBRACKET, Line: row + 1})
				index++
			} else if '0' <= line[index] && line[index] <= '9' {
				if line[index] == '0' {
					if index+1 < len(line)-1 && ('1' <= line[index+1] && line[index+1] <= '9') {
						return index, []token.Token{}, errors.New(fmt.Sprintf("Variable in Dfarg: %d find invalid token in \"%s\".", index, line))
					}
					tokens = append(tokens, token.Token{Content: "0", Kind: token.NUMBER, Line: row + 1})
					index++
					continue
				}

				start := index
				for index < len(line)-1 {
					if line[index] < '0' || '9' < line[index] {
						break
					}

					index++
				}
				tokens = append(tokens, token.Token{Content: line[start:index], Kind: token.NUMBER, Line: row + 1})
			} else if ('A' <= line[index] && line[index] <= 'Z') || ('a' <= line[index] && line[index] <= 'z') {
				start := index
				for index < len(line)-1 {
					if line[index] < 'A' || ('Z' < line[index] && line[index] < 'a') || 'z' < line[index] {
						break
					}

					index++
				}

				tokens = append(tokens, token.Token{Content: line[start:index], Kind: token.IDENTIFIER, Line: row + 1})
			} else if line[index] == '"' {
				var newToken token.Token
				var err error
				index, newToken, err = readString(index, line, row)
				if err != nil {
					return -1, nil, err
				}

				tokens = append(tokens, newToken)
			}
		}

		return index, []token.Token{}, errors.New(fmt.Sprintf("Variable in Dfarg: %d find invalid token in \"%s\".", index, line))
	} else {
		start := index
		for index < len(line)-1 {
			if line[index:index+2] == "{{" {
				break
			}
			index++
		}

		if index == len(line)-1 {
			return len(line), []token.Token{{Content: line[start:], Kind: token.DFARG, Line: row + 1}}, nil
		} else {
			return index, []token.Token{{Content: line[start:index], Kind: token.DFARG, Line: row + 1}}, nil
		}
	}
}

func readDfArgs(index int, line string, row int) (int, []token.Token, error) {
	var tokens []token.Token
	var newToken []token.Token
	var err error

	for line[index] == ' ' {
		index++
	}

	for index < len(line) {

		index, newToken, err = readDfArg(index, line, row)
		if err != nil {
			return index, tokens, err
		}

		tokens = append(tokens, newToken...)
	}

	return index, tokens, nil
}

func readString(index int, line string, row int) (int, token.Token, error) {
	start := index
	for index < len(line) {
		index++
		if line[index] == '"' {
			return index + 1, token.Token{Content: line[start+1 : index], Kind: token.STRING, Line: row + 1}, nil
		}
	}
	return index, token.Token{}, errors.New(fmt.Sprintf("tokenize error: found invalid token in line %d", row))
}

func readIdentifier(index int, line string, row int) (int, token.Token, error) {
	start := index
	if ('a' <= line[index] && line[index] <= 'z') || ('A' <= line[index] && line[index] <= 'Z') {
		for index < len(line) {
			index++
			if !isAlphabetOrNumber(line, index) {
				return index, token.Token{Content: line[start:index], Kind: token.IDENTIFIER, Line: row + 1}, nil
			}
		}
	}
	return index, token.Token{}, errors.New(fmt.Sprintf("tokenize error: found invalid token in line %d", row))
}
