package parser

import (
	"fmt"
	"errors"

	"dcc/types"
)

func program(tokens []types.Token, index int) error {
	var err error

	// { 関数インポート文 }
	for tokens[index].Kind == types.SIMPORT {
		index, err = importFunc(tokens, index)
		if err != nil {
			return err
		}
	}

	// { 関数 }
	for index < len(tokens) {
		if tokens[index].Kind != types.SIDENTIFIER {
			break
		}

		index, err = function(tokens, index)
		if err != nil {
			return err
		}
	}

	if index >= len(tokens) || tokens[index].Kind != types.SMAIN {
		return nil
	}
	
	// メイン部
	index, err = mainFunction(tokens, index)
	if err != nil {
		return err
	}

	return nil
}

// 関数インポート文
func importFunc(tokens []types.Token, index int) (int, error) {
	var err error
	// "import"
	if tokens[index].Kind != types.SIMPORT {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find 'import' in line %d", tokens[index].Line))
	}

	index++

	// 関数名
	index, err = functionName(tokens, index)
	if err != nil {
		return index, err
	}

	// "from"
	if tokens[index].Kind != types.SFROM {
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
func fileName(tokens []types.Token, index int) (int, error) {
	if tokens[index].Kind != types.SSTRING {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find an identifier in line %d", tokens[index].Line))
	}

	index++

	return index, nil
}

// 関数
func function(tokens []types.Token, index int) (int, error) {
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
func mainFunction(tokens []types.Token, index int) (int, error) {
	var err error

	// "main"
	if tokens[index].Kind != types.SMAIN {
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
func argumentDecralation(tokens []types.Token, index int) (int, error) {
	var err error

	// "("
	if tokens[index].Kind != types.SLPAREN {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find '(' in line %d", tokens[index].Line))
	}

	index++

	// 引数群
	if 	tokens[index].Kind == types.SIDENTIFIER {
		index, err = arguments(tokens, index)
		if err != nil {
			return index, err
		}
	}

	// ")"
	if tokens[index].Kind != types.SRPAREN {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find ')' in line %d", tokens[index].Line))
	}

	index++

	return index, nil
}

// 引数群
func arguments(tokens []types.Token, index int) (int, error) {
	var err error

	// 変数
	index, err = variable(tokens, index)
	if err != nil {
		return index, err
	}

	for ;; {
		// ","
		if tokens[index].Kind != types.SCOMMA {
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
func variable(tokens []types.Token, index int) (int, error) {
	var err error

	// 変数名
	index, err = variableName(tokens, index)
	if err != nil {
		return index, err
	}

	// "[]"
	if tokens[index].Kind == types.SARRANGE {
		index++
	}

	return index, nil
}

// 関数記述部
func functionDescription(tokens []types.Token, index int) (int, error) {
	var err error

	// "{"
	if tokens[index].Kind != types.SLBRACE {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find '{' in line %d", tokens[index].Line))
	}

	index++

	// 記述部
	index, err = description(tokens, index)
	if err != nil {
		return index, err
	}

	// "}"
	if tokens[index].Kind != types.SRBRACE {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find '}' in line %d", tokens[index].Line))
	}

	index++

	return index, nil
}

// 記述部
func description(tokens []types.Token, index int) (int, error) {
	var err error

	// 記述ブロック
	index, err = descriptionBlock(tokens, index)
	if err != nil {
		return index, err
	}

	for ;; {
		if tokens[index].Kind != types.SDFCOMMAND && tokens[index].Kind != types.SIDENTIFIER {
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
func descriptionBlock(tokens []types.Token, index int) (int, error) {
	var err error

	if tokens[index].Kind == types.SDFCOMMAND {
		// Dfile文
		index, err = dockerFile(tokens, index)
		if err != nil {
			return index, err
		}
	} else if tokens[index].Kind == types.SIDENTIFIER {
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
func dockerFile(tokens []types.Token, index int) (int, error) {
	var err error
	// Df命令
	if tokens[index].Kind != types.SDFCOMMAND {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find a Dockerfile comamnd in line %d", tokens[index].Line))
	}

	index++

	// Df引数部
	index, err = dfArgs(tokens, index)
	if err != nil {
		return index, err
	}

	return index, nil
}

// Df引数部
func dfArgs(tokens []types.Token, index int) (int, error) {
	var err error
	index, err = dfArg(tokens, index)
	if err != nil {
		return index, err
	}

	for ;; {
		if tokens[index].Kind == types.SDFARG || tokens[index].Kind == types.SASSIGNVARIABLE {
			index, err = dfArg(tokens, index)
			if err != nil {
				return index, err
			}
		} else {
			break
		}
	}

	return index, nil
}

// Df引数
func dfArg(tokens []types.Token, index int) (int, error) {
	if tokens[index].Kind != types.SDFARG && tokens[index].Kind != types.SASSIGNVARIABLE {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find Df argument in line %d", tokens[index].Line))
	}

	index++

	return index, nil
}

// 関数呼び出し文
func functionCall(tokens []types.Token, index int) (int, error) {
	var err error

	// 関数名
	index, err = functionName(tokens, index)
	if err != nil {
		return index, err
	}

	// "("
	if tokens[index].Kind != types.SLPAREN {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find '(' in line %d", tokens[index].Line))
	}

	index++

	// 文字列の並び
	if tokens[index].Kind == types.SSTRING {
		index, err = rowOfStrings(tokens, index)
		if err != nil {
			return index, err
		}
	}

	// ")"
	if tokens[index].Kind != types.SRPAREN {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find ')' in line %d", tokens[index].Line))
	}

	index++

	return index, nil
}

// 文字列の並び
func rowOfStrings(tokens []types.Token, index int) (int, error) {
	// 文字列
	if tokens[index].Kind != types.SSTRING {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find a string in line %d", tokens[index].Line))
	}

	index++

	for ;; {
		if tokens[index].Kind != types.SCOMMA {
			break
		}

		index++

		if tokens[index].Kind != types.SSTRING {
			return index, errors.New(fmt.Sprintf("syntax error: cannot find a string in line %d", tokens[index].Line))
		}

		index++
	}

	return index, nil
}

// 関数名
func functionName(tokens []types.Token, index int) (int, error) {
	if tokens[index].Kind != types.SIDENTIFIER {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find an identifier in line %d", tokens[index].Line))
	}

	index++

	return index, nil
}

// 変数名
func variableName(tokens []types.Token, index int) (int, error) {
	if tokens[index].Kind != types.SIDENTIFIER {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find an identifier in line %d", tokens[index].Line))
	}

	index++

	return index, nil
}