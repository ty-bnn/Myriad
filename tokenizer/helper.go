package tokenizer

import(
	"fmt"
	"errors"
	. "strings"
)

func isAlphabetOrNumber(line string, index int) bool {
	if (len(line) <= index || (line[index] < '0' || ('9' < line[index] && line[index] < 'A') || ('Z' < line[index] && line[index] < 'a') || ('z' < line[index]))) {
		return false
	}

	return true
}

func readReservedWords(index int, line string, row int) (int, Token, error) {
	if (index + 6 <= len(line) && line[index : index + 6] == "import" && !isAlphabetOrNumber(line, index + 6)) {
		return index + 6, Token{Content: "import", Kind: SIMPORT, Line: row + 1}, nil
	} else if (index + 4 <= len(line) && line[index : index + 4] == "from" && !isAlphabetOrNumber(line, index + 4)) {
		return index + 4, Token{Content: "from", Kind: SFROM, Line: row + 1}, nil
	} else if (index + 4 <= len(line) && line[index : index + 4] == "main" && !isAlphabetOrNumber(line, index + 4)) {
		return index + 4, Token{Content: "main", Kind: SMAIN, Line: row + 1}, nil
	} else if (index + 2 <= len(line) && line[index : index + 2] == "if" && !isAlphabetOrNumber(line, index + 2)) {
		return index + 2, Token{Content: "if", Kind: SIF, Line: row + 1}, nil
	} else if (index + 7 <= len(line) && line[index : index + 7] == "else if" && !isAlphabetOrNumber(line, index + 7)) {
		return index + 7, Token{Content: "else if", Kind: SELIF, Line: row + 1}, nil
	} else if (index + 4 <= len(line) && line[index : index + 4] == "else" && !isAlphabetOrNumber(line, index + 4)) {
		return index + 4, Token{Content: "else", Kind: SELSE, Line: row + 1}, nil
	}

	return index, Token{}, errors.New(fmt.Sprintf("tokenize error: found invalid character in line %d", row)) 
}

func readSymbols(index int, line string, row int) (int, Token, error) {
	if (line[index] == '(') {
		return index + 1, Token{Content: "(", Kind: SLPAREN, Line: row + 1}, nil
	} else if (line[index] == ')') {
		return index + 1, Token{Content: ")", Kind: SRPAREN, Line: row + 1}, nil
	} else if (line[index] == ',') {
		return index + 1, Token{Content: ",", Kind: SCOMMA, Line: row + 1}, nil
	} else if (index + 2 <= len(line) && line[index : index + 2] == "[]") {
		return index + 2, Token{Content: "[]", Kind: SARRANGE, Line: row + 1}, nil
	} else if (line[index] == '{') {
		return index + 1, Token{Content: "{", Kind: SLBRACE, Line: row + 1}, nil
	} else if (line[index] == '}') {
		return index + 1, Token{Content: "}", Kind: SRBRACE, Line: row + 1}, nil
	} else if (index + 2 <= len(line) && line[index : index + 2] == "==") {
		return index + 2, Token{Content: "==", Kind: SEQUAL, Line: row + 1}, nil
	} else if (index + 2 <= len(line) && line[index : index + 2] == "!=") {
		return index + 2, Token{Content: "!=", Kind: SNOTEQUAL, Line: row + 1}, nil
	} else if (index + 2 <= len(line) && line[index : index + 2] == ":=") {
		return index + 2, Token{Content: ":=", Kind: SASSIGN, Line: row + 1}, nil
	}

	return index, Token{}, errors.New(fmt.Sprintf("Symbols'index: %d find invalid character in \"%s\".", index, line))
}

func readDfCommands(index int, line string, row int) (int, Token, error) {
	/*
	ADD, ARG, CMD, COPY, ENTRYPOINT, ENV, EXPOSE, FROM,
	HEALTHCHECK, LABEL, MAINTAINER, ONBUILD, RUN, SHELL,
	STOPSIGNAL, USER, VOLUME, WORKDIR
	*/
	if HasPrefix(line[index:], "ADD") || HasPrefix(line[index:], "ARG") || HasPrefix(line[index:], "CMD") || HasPrefix(line[index:], "ENV") || HasPrefix(line[index:], "RUN") {
			return index + 3, Token{Content: line[index : index + 3], Kind: SDFCOMMAND, Line: row + 1}, nil
	}

	if HasPrefix(line[index:], "COPY") || HasPrefix(line[index:], "FROM") || HasPrefix(line[index:], "USER") {
		return index + 4, Token{Content: line[index : index + 4], Kind: SDFCOMMAND, Line: row + 1}, nil
	}

	if HasPrefix(line[index:], "LABEL") || HasPrefix(line[index:], "SHELL") {
		return index + 5, Token{Content: line[index : index + 5], Kind: SDFCOMMAND, Line: row + 1}, nil
	}

	if HasPrefix(line[index:], "") || HasPrefix(line[index:], "") {
		return index + 6, Token{Content: line[index : index + 6], Kind: SDFCOMMAND, Line: row + 1}, nil
	}

	if HasPrefix(line[index:], "ONBUILD") || HasPrefix(line[index:], "WORKDIR") {
		return index + 7, Token{Content: line[index : index + 7], Kind: SDFCOMMAND, Line: row + 1}, nil
	}

	if HasPrefix(line[index:], "ENTRYPOINT") || HasPrefix(line[index:], "MAINTAINER") || HasPrefix(line[index:], "STOPSIGNAL") {
		return index + 10, Token{Content: line[index : index + 10], Kind: SDFCOMMAND, Line: row + 1}, nil
	}

	if HasPrefix(line[index:], "HEALTHCHECK") {
		return index + 11, Token{Content: line[index : index + 10], Kind: SDFCOMMAND, Line: row + 1}, nil
	}

	return index, Token{}, errors.New(fmt.Sprintf("tokenize error: found invalid character in line %d", row))
}

func readDfArg(index int, line string, row int) (int, Token, error) {
	start := index

	if index + 2 < len(line) && line[index : index + 2] == "{{" {
		index += 2
		for index < len(line) - 1 {
			if line[index : index + 2] == "}}" {
				break
			}
			index++
		}
		if index != len(line) {
			return index + 2, Token{Content: line[start + 2: index], Kind: SASSIGNVARIABLE, Line: row + 1}, nil
		} else {
			return index, Token{}, errors.New(fmt.Sprintf("Variable in Dfarg: %d find invalid token in \"%s\".", index, line))
		}
	} else {
		for index < len(line) - 1 {
			if line[index : index + 2] == "{{" {
				break
			}
			index++
		}
		if index == len(line) - 1 {
			return len(line), Token{Content: line[start : len(line)], Kind: SDFARG, Line: row + 1}, nil
		} else {
			return index, Token{Content: line[start : index], Kind: SDFARG, Line: row + 1}, nil
		}
	}
}

func readDfArgs(index int, line string, row int) (int, []Token, error) {
	var tokens []Token
	var token Token
	var err error

	for line[index] == ' ' {
		index++
	}
	
	for index < len(line) {

		index, token, err = readDfArg(index, line, row)
		if err != nil {
			return index, tokens, err
		}

		tokens = append(tokens, token)
	}

	return index, tokens, nil
}

func readString(index int, line string, row int) (int, Token, error) {
	start := index
	for index < len(line) {
		index++
		if (line[index] == '"') {
			return index + 1, Token{Content: line[start+1 : index], Kind: SSTRING, Line: row + 1}, nil
		}
	}
	return index, Token{}, errors.New(fmt.Sprintf("tokenize error: found invalid token in line %d", row))
}

func readIdentifier(index int, line string, row int) (int, Token, error) {
	start := index
	if (('a' <= line[index] && line[index] <= 'z') || ('A' <= line[index] && line[index] <= 'Z')) {
		for index < len(line) {
			index++
			if (!isAlphabetOrNumber(line, index)) {
				return index, Token{Content: line[start : index], Kind: SIDENTIFIER, Line: row + 1}, nil
			}
		}
	}
	return index, Token{}, errors.New(fmt.Sprintf("tokenize error: found invalid token in line %d", row))
}
