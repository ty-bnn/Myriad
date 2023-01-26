package tokenizer

import (
	"regexp"
	// "fmt"

	"dcc/types"
)

func Tokenize(lines []string) ([]types.Token, error) {
	var tokens []types.Token
	reFunctionCall := regexp.MustCompile(`.*\((".*", )*(".*")?\)$`)
	reFunctionCallBeginLine := regexp.MustCompile(`.*\((".*", )*(".*",)?$`)
	reFunctionCallCenterLine := regexp.MustCompile(`\s*(".*", )*(".*",)$`)
	reFunctionCallEndLine := regexp.MustCompile(`\s*(".*", )*(".*")\)$`)
	reFunctionDeclareStart := regexp.MustCompile(`.*\(.*\) {$`) // ifも含まれる
	reFunctionDeclareEnd := regexp.MustCompile(`\s*}$`) // ifも含まれる
	reElIf := regexp.MustCompile(`\s*} else if \(.*\) {$`)
	reElse := regexp.MustCompile(`\s*} else {$`)
	reImport := regexp.MustCompile(`import.*from.*".*"`)

	for i := 0; i < len(lines); i++ {
		var token types.Token
		var err error 
		index := 0
		for index < len(lines[i]) {
			switch lines[i][index] {
				case ' ':
					index++
					continue
				/*
				Read symbols '(', ')', ',', "[]", '{', '}', "==", "!="
				*/
				case '(', ')', ',', '[', '{', '}', '=', '!':
					if reFunctionCall.MatchString(lines[i]) || reFunctionDeclareStart.MatchString(lines[i]) ||
						reFunctionDeclareEnd.MatchString(lines[i]) || reElIf.MatchString(lines[i]) ||
						reElse.MatchString(lines[i]) || reImport.MatchString(lines[i]) ||
						reFunctionCallBeginLine.MatchString(lines[i]) || reFunctionCallCenterLine.MatchString(lines[i]) ||
						reFunctionCallEndLine.MatchString(lines[i]) {
						index, token, err = readSymbols(index, lines[i], i)
					} else {
						var dfArgTokens []types.Token
						index, dfArgTokens, err = readDfArgsPerLine(index, lines[i], i)
						tokens = append(tokens, types.Token{Content: "    ", Kind: types.SDFARG})
						tokens = append(tokens, dfArgTokens...)
						token = types.Token{Content: "\n", Kind: types.SDFARG}
					}
				/*
				Read reserved words ("import", "from", "main", "if", "else if", "else")
				If not reserved words, read identifiers starts from 'i', 'f', 'm'.
				*/
				case 'e', 'f', 'i', 'm':
					if reFunctionCall.MatchString(lines[i]) || reFunctionDeclareStart.MatchString(lines[i]) ||
						reFunctionDeclareEnd.MatchString(lines[i]) || reElIf.MatchString(lines[i]) ||
						reElse.MatchString(lines[i]) || reImport.MatchString(lines[i]) ||
						reFunctionCallBeginLine.MatchString(lines[i]) || reFunctionCallCenterLine.MatchString(lines[i]) ||
						reFunctionCallEndLine.MatchString(lines[i]) {
						index, token, err = readReservedWords(index, lines[i], i)
						if err != nil {
							index, token, err = readIdentifier(index, lines[i], i)
						}
					} else {
						var dfArgTokens []types.Token
						index, dfArgTokens, err = readDfArgsPerLine(index, lines[i], i)
						tokens = append(tokens, types.Token{Content: "    ", Kind: types.SDFARG})
						tokens = append(tokens, dfArgTokens...)
						token = types.Token{Content: "\n", Kind: types.SDFARG}
					}
				/*
				Read Dfile commands.
				And Read Dfile arguments.
				*/
				case 'A', 'C', 'E', 'F', 'H', 'L', 'M', 'O', 'S', 'U', 'V', 'W', 'R':
					index, token, err = readDfCommands(index, lines[i], i)
					if err == nil {
						tokens = append(tokens, token)
					} else {
						tokens = append(tokens, types.Token{Content: "    ", Kind: types.SDFARG})
					}
					
					var dfArgTokens []types.Token
					index, dfArgTokens, err = readDfArgsPerLine(index, lines[i], i)
					tokens = append(tokens, dfArgTokens...)
					token = types.Token{Content: "\n", Kind: types.SDFARG}
				/*
				Read strings start from " and ends at ".
				*/
				case '"':
					index, token, err = readString(index, lines[i], i)
				/*
				Read identifiers.
				*/
				default:
					if reFunctionCall.MatchString(lines[i]) || reFunctionDeclareStart.MatchString(lines[i]) ||
						reFunctionDeclareEnd.MatchString(lines[i]) || reElIf.MatchString(lines[i]) ||
						reElse.MatchString(lines[i]) || reImport.MatchString(lines[i]) ||
						reFunctionCallBeginLine.MatchString(lines[i]) || reFunctionCallCenterLine.MatchString(lines[i]) ||
						reFunctionCallEndLine.MatchString(lines[i]) {
						index, token, err = readIdentifier(index, lines[i], i)
					} else {
						var dfArgTokens []types.Token
						index, dfArgTokens, err = readDfArgsPerLine(index, lines[i], i)
						tokens = append(tokens, types.Token{Content: "    ", Kind: types.SDFARG})
						tokens = append(tokens, dfArgTokens...)
						token = types.Token{Content: "\n", Kind: types.SDFARG}
					}
			}

			if err != nil {
				return nil, err
			}

			tokens = append(tokens, token)
		}
	}

	return tokens, nil
}