package tokenizer

import "github.com/ty-bnn/myriad/pkg/model/token"

func isLetter(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z')
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isWhiteSpace(ch byte) bool {
	return ch == ' ' || ch == '\t'
}

func isNewLine(ch byte) bool {
	return ch == '\n'
}

func getIdentKind(keyword string) token.TokenKind {
	kind, ok := token.ReservedKeywords[keyword]
	if !ok {
		return token.IDENTIFIER
	}

	return kind
}

func isDockerfileCommand(keyword string) bool {
	_, ok := token.DockerfileCommands[keyword]

	return ok
}

func (t *Tokenizer) nextTokenIs(word string) bool {
	return t.p+len(word)-1 < len(t.data) && t.data[t.p:t.p+len(word)] == word
}
