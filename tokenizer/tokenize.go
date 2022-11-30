package tokenizer

type TokenKind int

const (
	SIMPORT TokenKind = iota
	SFROM
	SMAIN
	SLPAREN
	SRPAREN
	SCOMMA
	SARRANGE
	SLBRACE
	SRBRACE
	SSTRING
	SDFCOMMAND
	SDFARG
	SIDENTIFIER
	SVARIABLE
)

type Token struct {
	Content string
	Kind TokenKind
	Line int
}

func Tokenize(lines []string) ([]Token, error) {
	var tokens []Token

	for line, lineStr := range lines {
		var token Token
		var err error 
		index := 0
		for index < len(lineStr) {
			switch lineStr[index] {
				case ' ':
					index++
					continue
				/*
				Read symbols '(', ')', ',', "[]"", '{', '}'
				*/
				case '(', ')', ',', '[', '{', '}':
					index, token, err = readSymbols(index, lineStr, line)
				/*
				Read reserved words ("import", "from", "main")
				If not reserved words, read identifiers starts from 'i', 'f', 'm'.
				*/
				case 'i', 'f', 'm':
					index, token, err = readReservedWords(index, lineStr, line)
					if err != nil {
						index, token, err = readIdentifier(index, lineStr, line)
					}				
				/*
				Read Dfile commands.
				And Read Dfile arguments.
				*/
				case 'A', 'C', 'E', 'F', 'H', 'L', 'M', 'O', 'R', 'S', 'U', 'V', 'W':
					index, token, err = readDfCommands(index, lineStr, line)
					if err == nil {
						for index < len(lineStr) {
							tokens = append(tokens, token)
							for index < len(lineStr) {
								if lineStr[index] != ' ' {
									break
								}
								index++
							}
							index, token, err = readDfArgs(index, lineStr, line)
							if err != nil {
								break
							}
						}
					}
				/*
				Read strings start from " and ends at ".
				*/
				case '"':
					index, token, err = readString(index, lineStr, line)
				/*
				Read identifiers.
				*/
				default:
					index, token, err = readIdentifier(index, lineStr, line)
			}

			if err != nil {
				return nil, err
			}

			tokens = append(tokens, token)
		}
	}

	return tokens, nil
}