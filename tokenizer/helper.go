package tokenizer

import(
	"fmt"
	"errors"
)

func readReservedWords(index int, lineStr string, line int) (int, Token, error) {
	if (lineStr[index : index + 6] == "import") {
		return index + 6, Token{"import", SIMPORT, line + 1}, nil 
	} else if (lineStr[index : index + 4] == "from") {
		return index + 4, Token{"from", SFROM, line + 1}, nil
	} else if (lineStr[index : index + 4] == "main") {
		return index + 4, Token{"main", SMAIN, line + 1}, nil
	}

	return index, Token{}, errors.New(fmt.Sprintf("ReservedWords'index: %d find invalid token in \"%s\".", index, lineStr)) 
}

func readSymbols(index int, lineStr string, line int) (int, Token, error) {
	if (lineStr[index] == '(') {
		return index + 1, Token{"(", SLPAREN, line + 1}, nil
	} else if (lineStr[index] == ')') {
		return index + 1, Token{")", SRPAREN, line + 1}, nil
	} else if (lineStr[index] == ',') {
		return index + 1, Token{",", SCOMMA, line + 1}, nil
	} else if (lineStr[index : index + 1] == "[]") {
		return index + 2, Token{"[]", SARRANGE, line + 1}, nil
	} else if (lineStr[index] == '{') {
		return index + 1, Token{"{", SLBRACE, line + 1}, nil
	} else if (lineStr[index] == '}') {
		return index + 1, Token{"}", SRBRACE, line + 1}, nil
	}

	return index, Token{}, errors.New(fmt.Sprintf("Symbols'index: %d find invalid token in \"%s\".", index, lineStr))
}

func readDfCommands(index int, lineStr string, line int) (int, Token, error) {
	/*
	ADD, ARG, CMD, COPY, ENTRYPOINT, ENV, EXPOSE, FROM,
	HEALTHCHECK, LABEL, MAINTAINER, ONBUILD, RUN, SHELL,
	STOPSIGNAL, USER, VOLUME, WORKDIR
	*/
	switch lineStr[index : index + 4] {
		case "ADD ", "ARG ", "CMD ", "ENV ", "RUN ":
			return index + 4, Token{lineStr[index : index + 3], SDFCOMMAND, line + 1}, nil
	}
	switch lineStr[index : index + 5] {
		case "COPY ", "FROM ", "USER ":
			return index + 5, Token{lineStr[index : index + 4], SDFCOMMAND, line + 1}, nil
	}
	switch lineStr[index : index + 6] {
		case "LABEL ", "SHELL ":
			return index + 6, Token{lineStr[index : index + 5], SDFCOMMAND, line + 1}, nil
	}
	switch lineStr[index : index + 7] {
		case "EXPOSE ", "VOLUME ":
			return index + 7, Token{lineStr[index : index + 6], SDFCOMMAND, line + 1}, nil
	}
	switch lineStr[index : index + 8] {
		case "ONBUILD ", "WORKDIR ":
			return index + 8, Token{lineStr[index : index + 7], SDFCOMMAND, line + 1}, nil
	}
	switch lineStr[index : index + 11] {
		case "ENTRYPOINT ", "MAINTAINER ", "STOPSIGNAL ":
			return index + 11, Token{lineStr[index : index + 10], SDFCOMMAND, line + 1}, nil
	}
	switch lineStr[index : index + 12] {
		case "HEALTHCHECK ":
			return index + 11, Token{lineStr[index : index + 10], SDFCOMMAND, line + 1}, nil
	}

	return index, Token{}, errors.New(fmt.Sprintf("DfCommand'index: %d find invalid token in \"%s\".", index, lineStr))
}

func readDfArgs(index int, lineStr string, line int) (int, Token, error) {
	start := index

	if lineStr[index : index + 2] == "${" {
		index += 2
		for index < len(lineStr) {
			if lineStr[index] == '}' {
				break
			}
			index++
		}
		if index != len(lineStr) {
			return index + 1, Token{lineStr[start + 2: index], SVARIABLE, line + 1}, nil
		} else {
			return index, Token{}, errors.New(fmt.Sprintf("Variable in Dfarg: %d find invalid token in \"%s\".", index, lineStr))
		}
	} else {
		for index < len(lineStr) - 1 {
			if lineStr[index : index + 2] == "${" {
				break
			}
			index++
		}
		if index == len(lineStr) - 1 {
			return len(lineStr), Token{lineStr[start : len(lineStr)], SDFARG, line + 1}, nil
		} else {
			return index, Token{lineStr[start : index], SDFARG, line + 1}, nil
		}
	}
}

func readString(index int, lineStr string, line int) (int, Token, error) {
	start := index
	for index < len(lineStr) {
		index++
		if (lineStr[index] == '"') {
			return index + 1, Token{lineStr[start+1 : index], SSTRING, line + 1}, nil
		}
	}
	return index, Token{}, errors.New(fmt.Sprintf("String's index: %d find invalid token in \"%s\".", index, lineStr))
}

func readIdentifier(index int, lineStr string, line int) (int, Token, error) {
	start := index
	if ('a' <= lineStr[index] && lineStr[index] <= 'z') {
		for index < len(lineStr) {
			index++
			if (lineStr[index] < 'A' || ('Z' < lineStr[index] && lineStr[index] < 'a') || 'z' < lineStr[index]) {
				return index, Token{lineStr[start : index], SIDENTIFIER, line + 1}, nil
			}
		}
	}
	return index, Token{}, errors.New(fmt.Sprintf("Others'index: %d find invalid token in \"%s\".", index, lineStr))
}