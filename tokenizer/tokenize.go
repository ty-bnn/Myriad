package tokenizer

import(
)

type Token struct {
	Content string
	Kind string
}

func Tokenize(lines []string) ([]Token, error) {
	var tokens []Token

	for _, line := range lines {
		var token Token
		var err error 
		i := 0
		for i < len(line) {
			switch line[i] {
				case ' ':
					i++
					continue
				/*
				Read symbols '(', ')', ',', "[]"", '{', '}'
				*/
				case '(', ')', ',', '[', '{', '}':
					i, token, err = readSymbols(i, line)
				/*
				Read reserved words ("import", "from", "main")
				If not reserved words, read identifiers starts from 'i', 'f', 'm'.
				*/
				case 'i', 'f', 'm':
					i, token, err = readReservedWords(i, line)
					if err != nil {
						i, token, err = readIdentifier(i, line)
					}				
				/*
				Read Dfile commands.
				And Read Dfile arguments.
				*/
				case 'A', 'C', 'E', 'F', 'H', 'L', 'M', 'O', 'R', 'S', 'U', 'V', 'W':
					i, token, err = readDfCommands(i, line)
					if err == nil {
						tokens = append(tokens, token)
						i, token, err = readDfArgs(i, line)
					}
				/*
				Read strings start from " and ends at ".
				*/
				case '"':
					i, token, err = readString(i, line)
				/*
				Read identifiers.
				*/
				default:
					i, token, err = readIdentifier(i, line)
			}

			if err != nil {
				return nil, err
			}

			tokens = append(tokens, token)
		}
	}

	return tokens, nil
}