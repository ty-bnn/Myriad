package parser

import (
	"fmt"
	"errors"

	"dcc/tokenizer"
)

func program(tokens []tokenizer.Token, index int) error {
	var err error

	// { 関数インポート文 }
	for tokens[index].Kind == tokenizer.SIMPORT {
		index, err = importFunc(tokens, index)
		if err != nil {
			return err
		}
	}

	// { 関数 }
	for tokens[index].Kind == tokenizer.SIDENTIFIER {
		index, err = function(tokens, index)
		if err != nil {
			return err
		}
	}

	// メイン部
	index, err = mainFunction(tokens, index)
	if err != nil {
		return err
	}

	return nil
}

// 関数インポート文
func importFunc(tokens []tokenizer.Token, index int) (int, error) {
	var err error
	// "import"
	if tokens[index].Kind != tokenizer.SIMPORT {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find 'import' in line %d", tokens[index].Line))
	}

	index++

	// 関数名
	index, err = functionName(tokens, index)
	if err != nil {
		return index, err
	}

	// "from"
	if tokens[index].Kind != tokenizer.SFROM {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find 'from' in line %d", tokens[index].Line))
	}

	index++

	// ファイル名
	index, err = fileName(tokens, index)
	if err != nil {
		return index, err
	}

	return index, nil
}

// ファイル名
func fileName(tokens []tokenizer.Token, index int) (int, error) {
	if tokens[index].Kind != tokenizer.SSTRING {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find an identifier in line %d", tokens[index].Line))
	}

	index++

	return index, nil
}

// 関数
func function(tokens []tokenizer.Token, index int) (int, error) {
	var err error

	// 関数名
	index, err = functionName(tokens, index)
	if err != nil {
		return index, err
	}

	// 引数宣言部
	index, err = argumentDecralation(tokens, index)
	if err != nil {
		return index, err
	}

	// 関数記述部
	index, err = functionDescription(tokens, index)
	if err != nil {
		return index, err
	}

	return index, nil
}

// メイン部
func mainFunction(tokens []tokenizer.Token, index int) (int, error) {
	var err error

	// "main"
	if tokens[index].Kind != tokenizer.SMAIN {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find 'main' to declare in line %d", tokens[index].Line))
	}

	index++

	// 引数宣言部
	index, err = argumentDecralation(tokens, index)
	if err != nil {
		return index, err
	}

	// 関数記述部
	index, err = functionDescription(tokens, index)
	if err != nil {
		return index, err
	}

	return index, nil
}

// 引数宣言部
func argumentDecralation(tokens []tokenizer.Token, index int) (int, error) {
	var err error

	// "("
	if tokens[index].Kind != tokenizer.SLPAREN {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find '(' in line %d", tokens[index].Line))
	}

	index++

	// 引数群
	if 	tokens[index].Kind == tokenizer.SIDENTIFIER {
		index, err = arguments(tokens, index)
		if err != nil {
			return index, err
		}
	}

	// ")"
	if tokens[index].Kind != tokenizer.SRPAREN {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find ')' in line %d", tokens[index].Line))
	}

	index++

	return index, nil
}

// 引数群
func arguments(tokens []tokenizer.Token, index int) (int, error) {
	var err error

	// 変数
	index, err = variable(tokens, index)
	if err != nil {
		return index, err
	}

	for ;; {
		// ","
		if tokens[index].Kind != tokenizer.SCOMMA {
			break
		}

		index++
		
		// 変数
		index, err = variable(tokens, index)
		if err != nil {
			return index, err
		}
	}

	return index, nil
}

// 変数
func variable(tokens []tokenizer.Token, index int) (int, error) {
	var err error

	// 変数名
	index, err = variableName(tokens, index)
	if err != nil {
		return index, err
	}

	// "[]"
	if tokens[index].Kind == tokenizer.SARRANGE {
		index++
	}

	return index, nil
}

// 関数記述部
func functionDescription(tokens []tokenizer.Token, index int) (int, error) {
	var err error

	// "{"
	if tokens[index].Kind != tokenizer.SLBRACE {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find '{' in line %d", tokens[index].Line))
	}

	index++

	// 記述部
	index, err = description(tokens, index)
	if err != nil {
		return index, err
	}

	// "}"
	if tokens[index].Kind != tokenizer.SRBRACE {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find '}' in line %d", tokens[index].Line))
	}

	index++

	return index, nil
}

// 記述部
func description(tokens []tokenizer.Token, index int) (int, error) {
	var err error

	// 記述ブロック
	index, err = descriptionBlock(tokens, index)
	if err != nil {
		return index, err
	}

	for ;; {
		if tokens[index].Kind != tokenizer.SDFCOMMAND && tokens[index].Kind != tokenizer.SIDENTIFIER {
			break
		}
			
		index, err = descriptionBlock(tokens, index)
		if err != nil {
			return index, err
		}
	}

	return index, nil
}

// 記述ブロック
func descriptionBlock(tokens []tokenizer.Token, index int) (int, error) {
	var err error

	if tokens[index].Kind == tokenizer.SDFCOMMAND {
		// Dfile文
		index, err = dockerFile(tokens, index)
		if err != nil {
			return index, err
		}
	} else if tokens[index].Kind == tokenizer.SIDENTIFIER {
		// 関数呼び出し文
		index, err = functionCall(tokens, index)
		if err != nil {
			return index, err
		}
	} else {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find a description block in line %d", tokens[index].Line))
	}

	return index, nil
}

// Dfile文
func dockerFile(tokens []tokenizer.Token, index int) (int, error) {
	// Df命令
	if tokens[index].Kind != tokenizer.SDFCOMMAND {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find a Dockerfile comamnd in line %d", tokens[index].Line))
	}

	index++

	// Df引数
	if tokens[index].Kind != tokenizer.SDFARG {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find a Dockerfile argument in line %d", tokens[index].Line))
	}

	index++

	return index, nil
}

// 関数呼び出し文
func functionCall(tokens []tokenizer.Token, index int) (int, error) {
	var err error

	// 関数名
	index, err = functionName(tokens, index)
	if err != nil {
		return index, err
	}

	// "("
	if tokens[index].Kind != tokenizer.SLPAREN {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find '(' in line %d", tokens[index].Line))
	}

	index++

	// 文字列の並び
	if tokens[index].Kind == tokenizer.SSTRING {
		index, err = rowOfStrings(tokens, index)
		if err != nil {
			return index, err
		}
	}

	// ")"
	if tokens[index].Kind != tokenizer.SRPAREN {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find ')' in line %d", tokens[index].Line))
	}

	index++

	return index, nil
}

// 文字列の並び
func rowOfStrings(tokens []tokenizer.Token, index int) (int, error) {
	// 文字列
	if tokens[index].Kind != tokenizer.SSTRING {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find a string in line %d", tokens[index].Line))
	}

	index++

	for ;; {
		if tokens[index].Kind != tokenizer.SCOMMA {
			break
		}

		index++

		if tokens[index].Kind != tokenizer.SSTRING {
			return index, errors.New(fmt.Sprintf("syntax error: cannot find a string in line %d", tokens[index].Line))
		}

		index++
	}

	return index, nil
}

// 関数名
func functionName(tokens []tokenizer.Token, index int) (int, error) {
	if tokens[index].Kind != tokenizer.SIDENTIFIER {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find an identifier in line %d", tokens[index].Line))
	}

	index++

	return index, nil
}

// 変数名
func variableName(tokens []tokenizer.Token, index int) (int, error) {
	if tokens[index].Kind != tokenizer.SIDENTIFIER {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find an identifier in line %d", tokens[index].Line))
	}

	index++

	return index, nil
}