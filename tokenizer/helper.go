package tokenizer

import(
	"fmt"
	"errors"

	"dcc/types"
)

func readReservedWords(index int, lineStr string, line int) (int, types.Token, error) {
	if (index + 6 <= len(lineStr) && lineStr[index : index + 6] == "import") {
		return index + 6, types.Token{Content: "import", Kind: types.SIMPORT, Line: line + 1}, nil 
	} else if (index + 4 <= len(lineStr) && lineStr[index : index + 4] == "from") {
		return index + 4, types.Token{Content: "from", Kind: types.SFROM, Line: line + 1}, nil
	} else if (index + 4 <= len(lineStr) && lineStr[index : index + 4] == "main") {
		return index + 4, types.Token{Content: "main", Kind: types.SMAIN, Line: line + 1}, nil
	} else if (index + 2 <= len(lineStr) && lineStr[index : index + 2] == "if") {
		return index + 2, types.Token{Content: "if", Kind: types.SIF, Line: line + 1}, nil
	} else if (index + 7 <= len(lineStr) && lineStr[index : index + 7] == "else if") {
		return index + 7, types.Token{Content: "else if", Kind: types.SELIF, Line: line + 1}, nil
	} else if (index + 4 <= len(lineStr) && lineStr[index : index + 4] == "else") {
		return index + 4, types.Token{Content: "else", Kind: types.SELSE, Line: line + 1}, nil
	}

	return index, types.Token{}, errors.New(fmt.Sprintf("ReservedWords'index: %d find invalid token in \"%s\".", index, lineStr)) 
}

func readSymbols(index int, lineStr string, line int) (int, types.Token, error) {
	if (lineStr[index] == '(') {
		return index + 1, types.Token{Content: "(", Kind: types.SLPAREN, Line: line + 1}, nil
	} else if (lineStr[index] == ')') {
		return index + 1, types.Token{Content: ")", Kind: types.SRPAREN, Line: line + 1}, nil
	} else if (lineStr[index] == ',') {
		return index + 1, types.Token{Content: ",", Kind: types.SCOMMA, Line: line + 1}, nil
	} else if (index + 2 <= len(lineStr) && lineStr[index : index + 2] == "[]") {
		return index + 2, types.Token{Content: "[]", Kind: types.SARRANGE, Line: line + 1}, nil
	} else if (lineStr[index] == '{') {
		return index + 1, types.Token{Content: "{", Kind: types.SLBRACE, Line: line + 1}, nil
	} else if (lineStr[index] == '}') {
		return index + 1, types.Token{Content: "}", Kind: types.SRBRACE, Line: line + 1}, nil
	} else if (index + 2 <= len(lineStr) && lineStr[index : index + 2] == "==") {
		return index + 2, types.Token{Content: "==", Kind: types.SEQUAL, Line: line + 1}, nil
	} else if (index + 2 <= len(lineStr) && lineStr[index : index + 2] == "!=") {
		return index + 2, types.Token{Content: "!=", Kind: types.SNOTEQUAL, Line: line + 1}, nil
	}

	return index, types.Token{}, errors.New(fmt.Sprintf("Symbols'index: %d find invalid token in \"%s\".", index, lineStr))
}

func readDfCommands(index int, lineStr string, line int) (int, types.Token, error) {
	/*
	ADD, ARG, CMD, COPY, ENTRYPOINT, ENV, EXPOSE, FROM,
	HEALTHCHECK, LABEL, MAINTAINER, ONBUILD, RUN, SHELL,
	STOPSIGNAL, USER, VOLUME, WORKDIR
	*/
	switch lineStr[index : index + 4] {
		case "ADD ", "ARG ", "CMD ", "ENV ", "RUN ":
			return index + 4, types.Token{Content: lineStr[index : index + 3], Kind: types.SDFCOMMAND, Line: line + 1}, nil
	}
	switch lineStr[index : index + 5] {
		case "COPY ", "FROM ", "USER ":
			return index + 5, types.Token{Content: lineStr[index : index + 4], Kind: types.SDFCOMMAND, Line: line + 1}, nil
	}
	switch lineStr[index : index + 6] {
		case "LABEL ", "SHELL ":
			return index + 6, types.Token{Content: lineStr[index : index + 5], Kind: types.SDFCOMMAND, Line: line + 1}, nil
	}
	switch lineStr[index : index + 7] {
		case "EXPOSE ", "VOLUME ":
			return index + 7, types.Token{Content: lineStr[index : index + 6], Kind: types.SDFCOMMAND, Line: line + 1}, nil
	}
	switch lineStr[index : index + 8] {
		case "ONBUILD ", "WORKDIR ":
			return index + 8, types.Token{Content: lineStr[index : index + 7], Kind: types.SDFCOMMAND, Line: line + 1}, nil
	}
	switch lineStr[index : index + 11] {
		case "ENTRYPOINT ", "MAINTAINER ", "STOPSIGNAL ":
			return index + 11, types.Token{Content: lineStr[index : index + 10], Kind: types.SDFCOMMAND, Line: line + 1}, nil
	}
	switch lineStr[index : index + 12] {
		case "HEALTHCHECK ":
			return index + 11, types.Token{Content: lineStr[index : index + 10], Kind: types.SDFCOMMAND, Line: line + 1}, nil
	}

	return index, types.Token{}, errors.New(fmt.Sprintf("DfCommand'index: %d find invalid token in \"%s\".", index, lineStr))
}

func readDfArgs(index int, lineStr string, line int) (int, types.Token, error) {
	start := index

	if index + 2 < len(lineStr) && lineStr[index : index + 2] == "${" {
		index += 2
		for index < len(lineStr) {
			if lineStr[index] == '}' {
				break
			}
			index++
		}
		if index != len(lineStr) {
			return index + 1, types.Token{Content: lineStr[start + 2: index], Kind: types.SASSIGNVARIABLE, Line: line + 1}, nil
		} else {
			return index, types.Token{}, errors.New(fmt.Sprintf("Variable in Dfarg: %d find invalid token in \"%s\".", index, lineStr))
		}
	} else {
		for index < len(lineStr) - 1 {
			if lineStr[index : index + 2] == "${" {
				break
			}
			index++
		}
		if index == len(lineStr) - 1 {
			return len(lineStr), types.Token{Content: lineStr[start : len(lineStr)], Kind: types.SDFARG, Line: line + 1}, nil
		} else {
			return index, types.Token{Content: lineStr[start : index], Kind: types.SDFARG, Line: line + 1}, nil
		}
	}
}

func readDfArgsPerLine(index int, lineStr string, line int) (int, []types.Token, error) {
	var tokens []types.Token
	var token types.Token
	var err error

	for index < len(lineStr) {
		for lineStr[index] == ' ' {
			index++
		}

		index, token, err = readDfArgs(index, lineStr, line)
		if err != nil {
			return index, tokens, err
		}

		tokens = append(tokens, token)
	}

	return index, tokens, nil
}

func readString(index int, lineStr string, line int) (int, types.Token, error) {
	start := index
	for index < len(lineStr) {
		index++
		if (lineStr[index] == '"') {
			return index + 1, types.Token{Content: lineStr[start+1 : index], Kind: types.SSTRING, Line: line + 1}, nil
		}
	}
	return index, types.Token{}, errors.New(fmt.Sprintf("String's index: %d find invalid token in \"%s\".", index, lineStr))
}

func readIdentifier(index int, lineStr string, line int) (int, types.Token, error) {
	start := index
	if ('a' <= lineStr[index] && lineStr[index] <= 'z') {
		for index < len(lineStr) {
			index++
			if (lineStr[index] < 'A' || ('Z' < lineStr[index] && lineStr[index] < 'a') || 'z' < lineStr[index]) {
				return index, types.Token{Content: lineStr[start : index], Kind: types.SIDENTIFIER, Line: line + 1}, nil
			}
		}
	}
	return index, types.Token{}, errors.New(fmt.Sprintf("Others'index: %d find invalid token in \"%s\".", index, lineStr))
}