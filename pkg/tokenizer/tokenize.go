package tokenizer

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ty-bnn/myriad/pkg/model/token"
)

func (t *Tokenizer) Tokenize() error {
	fmt.Println("Tokenizing...")

	for t.p < len(t.data) {
		if !t.isInDfBlock {
			tok, err := t.TokenizeMyriad()
			if err != nil {
				return err
			}
			if tok != (token.Token{}) {
				t.Tokens = append(t.Tokens, tok)
			}
		} else {
			tok, err := t.TokenizeDockerfile()
			if err != nil {
				return err
			}
			if tok != (token.Token{}) {
				t.Tokens = append(t.Tokens, tok)
			}
		}
	}

	fmt.Println("Tokenize Done.")

	return nil
}

func (t *Tokenizer) TokenizeMyriad() (token.Token, error) {
	switch t.data[t.p] {
	case '(':
		t.p++
		return token.Token{Kind: token.LPAREN, Content: "("}, nil
	case ')':
		t.p++
		return token.Token{Kind: token.RPAREN, Content: ")"}, nil
	case ',':
		t.p++
		return token.Token{Kind: token.COMMA, Content: ","}, nil
	case '[':
		t.p++
		return token.Token{Kind: token.LBRACKET, Content: "["}, nil
	case ']':
		t.p++
		return token.Token{Kind: token.RBRACKET, Content: "]"}, nil
	case '{':
		if t.p+2 < len(t.data) && t.data[t.p:t.p+3] == "{{-" {
			t.p += 3
			t.isInDfBlock = true
			return token.Token{Kind: token.DFBEGIN, Content: "{{-"}, nil
		} else {
			t.p++
			return token.Token{Kind: token.LBRACE, Content: "{"}, nil
		}
	case '}':
		if t.p+1 < len(t.data) && t.data[t.p:t.p+2] == "}}" {
			t.p += 2
			t.isInDfBlock = true
			return token.Token{Kind: token.RDOUBLEBRA, Content: "}}"}, nil
		} else {
			t.p++
			return token.Token{Kind: token.RBRACE, Content: "}"}, nil
		}
	case '.':
		t.p++
		return token.Token{Kind: token.DOT, Content: "."}, nil
	case ':':
		if t.p+1 < len(t.data) && t.data[t.p:t.p+2] == ":=" {
			t.p += 2
			return token.Token{Kind: token.DEFINE, Content: ":="}, nil
		} else {
			return token.Token{}, errors.New(fmt.Sprintf("tokenize error: invalid token ':'"))
		}
	case '=':
		if t.p+1 < len(t.data) && t.data[t.p:t.p+2] == "==" {
			t.p += 2
			return token.Token{Kind: token.EQUAL, Content: "=="}, nil
		} else {
			t.p++
			return token.Token{Kind: token.ASSIGN, Content: "="}, nil
		}
	case '!':
		if t.p+1 < len(t.data) && t.data[t.p:t.p+2] == "!=" {
			t.p += 2
			return token.Token{Kind: token.NOTEQUAL, Content: "!="}, nil
		} else {
			t.p++
			return token.Token{Kind: token.NOT, Content: "!"}, nil
		}
	case '&':
		if t.p+1 < len(t.data) && t.data[t.p:t.p+2] == "&&" {
			t.p += 2
			return token.Token{Kind: token.AND, Content: "&&"}, nil
		} else {
			return token.Token{}, errors.New(fmt.Sprintf("tokenize error: invalid token '&'"))
		}
	case '|':
		if t.p+1 < len(t.data) && t.data[t.p:t.p+2] == "||" {
			t.p += 2
			return token.Token{Kind: token.OR, Content: "||"}, nil
		} else {
			return token.Token{}, errors.New(fmt.Sprintf("tokenize error: invalid token '|'"))
		}
	case '"':
		// TODO: パース時に判断したい
		t.p++
		start := t.p
		for t.p < len(t.data) && t.data[t.p] != '"' {
			t.p++
			if t.p == len(t.data) {
				return token.Token{}, errors.New(fmt.Sprintf("tokenize error: cannot find '\"'"))
			}
		}
		content := t.data[start:t.p]
		t.p++
		return token.Token{Kind: token.STRING, Content: content}, nil
	case '+':
		t.p++
		return token.Token{Kind: token.PLUS, Content: "+"}, nil
	case '<':
		if t.p+1 < len(t.data) && t.data[t.p:t.p+2] == "<<" {
			t.p += 2
			return token.Token{Kind: token.DOUBLELESS, Content: "<<"}, nil
		}
	default:
		if isWhiteSpace(t.data[t.p]) || isNewLine(t.data[t.p]) {
			t.p++
			return token.Token{}, nil
		} else if isLetter(t.data[t.p]) {
			start := t.p
			for t.p < len(t.data) && (isLetter(t.data[t.p]) || isDigit(t.data[t.p])) {
				t.p++
			}
			kind := getIdentKind(t.data[start:t.p])
			return token.Token{Kind: kind, Content: t.data[start:t.p]}, nil
		} else if isDigit(t.data[t.p]) {
			start := t.p
			for t.p < len(t.data) && isDigit(t.data[t.p]) {
				t.p++
			}
			return token.Token{Kind: token.NUMBER, Content: t.data[start:t.p]}, nil
		} else {
			return token.Token{}, errors.New(fmt.Sprintf("tokenize error: invalid token %b", t.data[t.p]))
		}
	}
	return token.Token{}, errors.New(fmt.Sprintf("tokenize error: invalid token"))
}

func (t *Tokenizer) TokenizeDockerfile() (token.Token, error) {
	if isWhiteSpace(t.data[t.p]) || isNewLine(t.data[t.p]) {
		t.p++
		return token.Token{}, nil
	}

	if t.p+2 < len(t.data) && t.data[t.p:t.p+3] == "-}}" {
		t.p += 3
		t.isInDfBlock = false
		return token.Token{Kind: token.DFEND, Content: "-}}"}, nil
	}

	if !t.isInCommand {
		start := t.p
		for t.p < len(t.data) && !isWhiteSpace(t.data[t.p]) && !isNewLine(t.data[t.p]) {
			t.p++
		}
		if isDockerfileCommand(t.data[start:t.p]) {
			t.commandPtr = t.data[start:t.p]
			t.isInCommand = true
			return token.Token{Kind: token.DFCOMMAND, Content: t.data[start:t.p]}, nil
		}

		return token.Token{}, errors.New(fmt.Sprintf("tokenize error: invalid token %s", t.data[start:t.p]))
	}

	start := t.p
	for t.p < len(t.data) {
		if isNewLine(t.data[t.p]) {
			// 末尾が'\'で終わっているか確認
			trimmed := strings.TrimSpace(t.data[start:t.p])
			if 0 < len(trimmed) && trimmed[len(trimmed)-1] == '\\' {
				t.isInCommand = true
			} else {
				t.isInCommand = false
			}
			t.p++
			return token.Token{Kind: token.DFARG, Content: t.data[start:t.p]}, nil
		} else if t.p+2 < len(t.data) && t.data[t.p:t.p+3] == "-}}" {
			// 末尾が'\'で終わっているか確認
			trimmed := strings.TrimSpace(t.data[start:t.p])
			if 0 < len(trimmed) && trimmed[len(trimmed)-1] == '\\' {
				t.isInCommand = true
			} else {
				t.isInCommand = false
			}
			return token.Token{Kind: token.DFARG, Content: strings.TrimRight(t.data[start:t.p], " ") + "\n"}, nil
		} else if t.p+1 < len(t.data) && t.data[t.p:t.p+2] == "{{" {
			t.p += 2
			t.isInDfBlock = false
			return token.Token{Kind: token.LDOUBLEBRA, Content: "{{"}, nil
		}
		t.p++
	}

	return token.Token{Kind: token.DFARG, Content: t.data[start:t.p]}, nil
}
