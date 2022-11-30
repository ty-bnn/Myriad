package compiler

import (
	"fmt"
	"errors"

	"dcc/tokenizer"
)

type CodeKind int

const (
	ROW CodeKind = iota
	VAR
)

type Variable struct {
	name string
}

type Function struct {
	codes []string
	argTable []Variable
}

type Code struct {
	code string
	aa CodeKind
}

var functionCodeMap map[string][]string
var functionPointer string

func program(tokens []tokenizer.Token, index int) error {
	functionCodeMap = map[string][]string{}

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
	index++

	// 関数名
	index, err = functionName(tokens, index)
	if err != nil {
		return index, err
	}

	// "from"
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
	index++

	return index, nil
}

// 関数
func function(tokens []tokenizer.Token, index int) (int, error) {
	var err error

	// 関数名
	functionPointer = tokens[index].Content

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
	index++

	// 引数群
	if 	tokens[index].Kind == tokenizer.SIDENTIFIER {
		index, err = arguments(tokens, index)
		if err != nil {
			return index, err
		}
	}

	// ")"
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
	index++

	// 記述部
	index, err = description(tokens, index)
	if err != nil {
		return index, err
	}

	// "}"
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
	}

	return index, nil
}

// Dfile文
func dockerFile(tokens []tokenizer.Token, index int) (int, error) {
	// Df命令
	code := tokens[index].Content + " "
	index++

	// Df引数
	code += tokens[index].Content + "\n"
	index++

	if functionPointer == "main" {
		dfCodes = append(dfCodes, code)
	} else {
		functionCodeMap[functionPointer] = append(functionCodeMap[functionPointer], code)
	}

	return index, nil
}

// 関数呼び出し文
func functionCall(tokens []tokenizer.Token, index int) (int, error) {
	var err error
	functionCallName := tokens[index].Content

	// 関数名
	if _, ok := functionCodeMap[functionCallName]; !ok {
		return index, errors.New(fmt.Sprintf("semantic error: cannot find '%s' in line %d", tokens[index].Content, tokens[index].Line))
	}
	index, err = functionName(tokens, index)
	if err != nil {
		return index, err
	}

	// "("
	index++

	// 文字列の並び
	if tokens[index].Kind == tokenizer.SSTRING {
		index, err = rowOfStrings(tokens, index)
		if err != nil {
			return index, err
		}
	}

	for _, code := range functionCodeMap[functionCallName] {
		if functionCallName == "main" {
			dfCodes = append(dfCodes, code)
		} else {
			functionCodeMap[functionPointer] = append(functionCodeMap[functionPointer], code)
		}
	}

	// ")"
	index++

	return index, nil
}

// 文字列の並び
func rowOfStrings(tokens []tokenizer.Token, index int) (int, error) {
	// 文字列
	index++

	for ;; {
		if tokens[index].Kind != tokenizer.SCOMMA {
			break
		}

		index++
		index++
	}

	return index, nil
}

// 関数名
func functionName(tokens []tokenizer.Token, index int) (int, error) {
	// 名前
	index++

	return index, nil
}

// 変数名
func variableName(tokens []tokenizer.Token, index int) (int, error) {
	// 名前
	index++

	return index, nil
}