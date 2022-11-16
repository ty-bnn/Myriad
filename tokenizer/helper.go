package tokenizer

import(
	"fmt"
	"errors"
)

func readReservedWords(index int, line string) (int, Token, error) {
	if (line[index : index + 6] == "import") {
		return index + 6, Token{"import", "SIMPORT"}, nil 
	} else if (line[index : index + 4] == "from") {
		return index + 4, Token{"from", "SFROM"}, nil
	} else if (line[index : index + 4] == "main") {
		return index + 4, Token{"main", "SMAIN"}, nil
	}

	return index, Token{}, errors.New(fmt.Sprintf("ReservedWords'index: %d find invalid token in \"%s\".", index, line)) 
}

func readSymbols(index int, line string) (int, Token, error) {
	if (line[index] == '(') {
		return index + 1, Token{"(", "SLPAREN"}, nil
	} else if (line[index] == ')') {
		return index + 1, Token{")", "SRPAREN"}, nil
	} else if (line[index] == ',') {
		return index + 1, Token{",", "SCOMMA"}, nil
	} else if (line[index : index + 1] == "[]") {
		return index + 2, Token{"[]", "SARRANGE"}, nil
	} else if (line[index] == '{') {
		return index + 1, Token{"{", "SLBRACE"}, nil
	} else if (line[index] == '}') {
		return index + 1, Token{"}", "SRBRACE"}, nil
	}

	return index, Token{}, errors.New(fmt.Sprintf("Symbols'index: %d find invalid token in \"%s\".", index, line))
}

func readDfCommands(index int, line string) (int, Token, error) {
	/*
	ADD, ARG, CMD, COPY, ENTRYPOINT, ENV, EXPOSE, FROM,
	HEALTHCHECK, LABEL, MAINTAINER, ONBUILD, RUN, SHELL,
	STOPSIGNAL, USER, VOLUME, WORKDIR
	*/
	switch line[index : index + 4] {
		case "ADD ", "ARG ", "CMD ", "ENV ", "RUN ":
			return index + 4, Token{line[index : index + 3], "SDFCMMAND"}, nil
	}
	switch line[index : index + 5] {
		case "COPY ", "FROM ", "USER ":
			return index + 5, Token{line[index : index + 4], "SDFCOMMAND"}, nil
	}
	switch line[index : index + 6] {
		case "LABEL ", "SHELL ":
			return index + 6, Token{line[index : index + 5], "SDFCOMMAND"}, nil
	}
	switch line[index : index + 7] {
		case "EXPOSE ", "VOLUME ":
			return index + 7, Token{line[index : index + 6], "SDFCOMMAND"}, nil
	}
	switch line[index : index + 8] {
		case "ONBUILD ", "WORKDIR ":
			return index + 8, Token{line[index : index + 7], "SDFCOMMAND"}, nil
	}
	switch line[index : index + 11] {
		case "ENTRYPOINT ", "MAINTAINER ", "STOPSIGNAL ":
			return index + 11, Token{line[index : index + 10], "SDFCOMMAND"}, nil
	}
	switch line[index : index + 12] {
		case "HEALTHCHECK ":
			return index + 11, Token{line[index : index + 10], "SDFCOMMAND"}, nil
	}

	return index, Token{}, errors.New(fmt.Sprintf("DfCommand'index: %d find invalid token in \"%s\".", index, line))
}

func readDfArgs(index int, line string) (int, Token, error) {
	for index < len(line) {
		if line[index] == ' ' {
			index++
			continue
		} else {
			break
		}
	}
	if index != len(line) {
		return len(line), Token{line[index : len(line)], "SDFARG"}, nil
	}
	
	return index, Token{}, errors.New(fmt.Sprintf("DfArg'index: %d find invalid token in \"%s\".", index, line))
}

func readString(index int, line string) (int, Token, error) {
	start := index
	for index < len(line) {
		index++
		if (line[index] == '"') {
			return index + 1, Token{line[start+1 : index], "SSTRING"}, nil
		}
	}
	return index, Token{}, errors.New(fmt.Sprintf("String'index: %d find invalid token in \"%s\".", index, line))
}

func readIdentifier(index int, line string) (int, Token, error) {
	start := index
	if ('a' <= line[index] && line[index] <= 'z') {
		for index < len(line) {
			index++
			if (line[index] < 'A' || ('Z' < line[index] && line[index] < 'a') || 'z' < line[index]) {
				return index, Token{line[start : index], "SIDENTIFIER"}, nil
			}
		}
	}
	return index, Token{}, errors.New(fmt.Sprintf("Others'index: %d find invalid token in \"%s\".", index, line))
}