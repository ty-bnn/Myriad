package compiler

import (
	"fmt"
	"os"
	"errors"

	"dcc/types"
	"dcc/tokenizer"
	"dcc/parser"
	"dcc/others"
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
func fileName(tokens []types.Token, index int) (int, error) {
	filePath := tokens[index].Content
	lines, err := others.ReadLinesFromFile(filePath)
	if err != nil {
		return index, err
	}

	newTokens, err := tokenizer.Tokenize(lines)
	if err != nil {
		return index, err
	}

	err = parser.Parse(newTokens)
	if err != nil {
		return index, err
	}

	err = Compile(newTokens, functionArgMap, functionCodeMap)
	if err != nil {
		return index, err
	}

	index++

	return index, nil
}

// 関数
func function(tokens []types.Token, index int) (int, error) {
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
func mainFunction(tokens []types.Token, index int) (int, error) {
	var err error

	functionPointer = "main"

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
func argumentDecralation(tokens []types.Token, index int) (int, error) {
	var err error

	// "("
	index++

	// 引数群
	if 	tokens[index].Kind == types.SIDENTIFIER {
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
func arguments(tokens []types.Token, index int) (int, error) {
	var err error
	argIndex := 2

	// 変数
	index, err = variable(tokens, index, argIndex)
	if err != nil {
		return index, err
	}

	for ;; {
		argIndex++
		// ","
		if tokens[index].Kind != types.SCOMMA {
			break
		}

		index++
		
		// 変数
		index, err = variable(tokens, index, argIndex)
		if err != nil {
			return index, err
		}
	}

	return index, nil
}

// 変数
func variable(tokens []types.Token, index int, argIndex int) (int, error) {
	var err error

	// 変数名
	name := tokens[index].Content
	var argument types.Argument
	if functionPointer == "main" {
		if argIndex < len(os.Args){
			argument = types.Argument{Name: name, Value: os.Args[argIndex], Kind: types.STRING}
		} else {
			return index, errors.New(fmt.Sprintf("semantic error: not enough argument value in line %d", tokens[index].Line))
		}
	} else {
		argument = types.Argument{Name: name, Value: "", Kind: types.STRING}
	}

	index, err = variableName(tokens, index)
	if err != nil {
		return index, err
	}

	// "[]"
	if tokens[index].Kind == types.SARRANGE {
		argument.Kind = types.ARRAY
		index++
	}

	if argumentExist(functionPointer, name) {
		return index, errors.New(fmt.Sprintf("semantic error: %s is already defined in line %d", name, tokens[index].Line))
	}

	(*functionArgMap)[functionPointer] = append((*functionArgMap)[functionPointer], argument)

	return index, nil
}

// 関数記述部
func functionDescription(tokens []types.Token, index int) (int, error) {
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
	}

	return index, nil
}

// Dfile文
func dockerFile(tokens []types.Token, index int) (int, error) {
	var err error
	// Df命令
	(*functionCodeMap)[functionPointer] = append((*functionCodeMap)[functionPointer], types.Code{Code: tokens[index].Content, Kind: types.ROW})
	(*functionCodeMap)[functionPointer] = append((*functionCodeMap)[functionPointer], types.Code{Code: " ", Kind: types.ROW})
	index++

	// Df引数
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
	content := tokens[index].Content

	var code types.Code
	if tokens[index].Kind == types.SDFARG {
		code = types.Code{Code: content, Kind: types.ROW}
	} else if tokens[index].Kind == types.SASSIGNVARIABLE && functionPointer == "main" {
		code = types.Code{Code: getArgumentValue("main", content)}
	} else {
		code = types.Code{Code: content, Kind: types.VAR}
	}
	
	(*functionCodeMap)[functionPointer] = append((*functionCodeMap)[functionPointer], code)
	index++

	return index, nil
}

// 関数呼び出し文
func functionCall(tokens []types.Token, index int) (int, error) {
	var err error
	functionCallName := tokens[index].Content

	// 関数名
	if _, ok := (*functionCodeMap)[functionCallName]; !ok {
		return index, errors.New(fmt.Sprintf("semantic error: function %s is not defined in line %d", tokens[index].Content, tokens[index].Line))
	}

	index, err = functionName(tokens, index)
	if err != nil {
		return index, err
	}

	// "("
	index++

	// 文字列の並び
	index, err = rowOfStrings(tokens, index)
	if err != nil {
		return index, err
	}

	for _, code := range (*functionCodeMap)[functionCallName] {
		if code.Kind == types.ROW {
			(*functionCodeMap)[functionPointer] = append((*functionCodeMap)[functionPointer], code)
		} else {
			(*functionCodeMap)[functionPointer] = append((*functionCodeMap)[functionPointer], types.Code{Code: getArgumentValue(functionCallName, code.Code)})
		}
	}

	// ")"
	index++

	return index, nil
}

// 文字列の並び
func rowOfStrings(tokens []types.Token, index int) (int, error) {
	functionCallName := tokens[index - 2].Content

	// 文字列
	for i, _ := range (*functionArgMap)[functionCallName] {
		if tokens[index].Kind != types.SSTRING {
			return index, errors.New(fmt.Sprintf("semantic error: not enough arguments in line %d", tokens[index].Line))
		}

		(*functionArgMap)[functionCallName][i].Value = tokens[index].Content
		index++

		if i == len((*functionArgMap)[functionCallName]) - 1 {
			break
		}

		index++
	}

	if tokens[index].Kind == types.SCOMMA {
		return index, errors.New(fmt.Sprintf("semantic error: too many arguments in line %d", tokens[index].Line))
	}

	return index, nil
}

// 関数名
func functionName(tokens []types.Token, index int) (int, error) {
	// 名前
	index++

	return index, nil
}

// 変数名
func variableName(tokens []types.Token, index int) (int, error) {
	// 名前
	index++

	return index, nil
}