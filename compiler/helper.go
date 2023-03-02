package compiler

import (
	"fmt"
	"errors"

	"dcc/types"
	"dcc/tokenizer"
	"dcc/parser"
	"dcc/others"
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
	for index < len(tokens) {
		if tokens[index].Kind != tokenizer.SIDENTIFIER {
			break
		}

		index, err = function(tokens, index)
		if err != nil {
			return err
		}
	}

	if index >= len(tokens) || tokens[index].Kind != tokenizer.SMAIN {
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
	filePath := tokens[index].Content

	if isCompiled(filePath) {
		return index + 1, nil
	}

	readFiles = append(readFiles, filePath)

	lines, err := others.ReadLinesFromFile(filePath)
	if err != nil {
		return index, err
	}

	newTokens, err := tokenizer.Tokenize(lines)
	if err != nil {
		return index, err
	}

	// For debug.
	// for _, token := range newTokens {
	// 	fmt.Printf("%30s\t%10d\n", token.Content, token.Kind)
	// }

	err = parser.Parse(newTokens)
	if err != nil {
		return index, err
	}

	err = program(newTokens, 0)
	if err != nil {
		return index, err
	}

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
	argIndex := 2

	// 変数
	index, err = variable(tokens, index, argIndex)
	if err != nil {
		return index, err
	}

	for ;; {
		argIndex++
		// ","
		if tokens[index].Kind != tokenizer.SCOMMA {
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
func variable(tokens []tokenizer.Token, index int, argIndex int) (int, error) {
	var err error

	// 変数名
	name := tokens[index].Content
	argument := types.Argument{Name: name, Kind: types.STRING}

	index, err = variableName(tokens, index)
	if err != nil {
		return index, err
	}

	// "[]"
	if tokens[index].Kind == tokenizer.SARRANGE {
		argument.Kind = types.ARRAY
		index++
	}

	functionArgMap[functionPointer] = append(functionArgMap[functionPointer], argument)

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
		if tokens[index].Kind != tokenizer.SDFCOMMAND && tokens[index].Kind != tokenizer.SDFARG && tokens[index].Kind != tokenizer.SIDENTIFIER && tokens[index].Kind != tokenizer.SIF {
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

	if tokens[index].Kind == tokenizer.SDFCOMMAND || tokens[index].Kind == tokenizer.SDFARG {
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
	} else if tokens[index].Kind == tokenizer.SIF {
		// ifブロック
		index, err = ifBlock(tokens, index)
		if err != nil {
			return index, err
		}
	}

	return index, nil
}

// Dfile文
func dockerFile(tokens []tokenizer.Token, index int) (int, error) {
	var err error
	if tokens[index].Kind == tokenizer.SDFCOMMAND {
		// Df命令
		functionInterCodeMap[functionPointer] = append(functionInterCodeMap[functionPointer], types.InterCode{Content: tokens[index].Content, Kind: types.ROW})
		functionInterCodeMap[functionPointer] = append(functionInterCodeMap[functionPointer], types.InterCode{Content: " ", Kind: types.ROW})
		index++
	}

	// Df引数
	index, err = dfArgs(tokens, index)
	if err != nil {
		return index, err
	}

	return index, nil
}

// Df引数部
func dfArgs(tokens []tokenizer.Token, index int) (int, error) {
	var err error
	index, err = dfArg(tokens, index)
	if err != nil {
		return index, err
	}

	for ;; {
		if tokens[index].Kind == tokenizer.SDFARG || tokens[index].Kind == tokenizer.SASSIGNVARIABLE {
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
func dfArg(tokens []tokenizer.Token, index int) (int, error) {
	content := tokens[index].Content

	var code types.InterCode
	if tokens[index].Kind == tokenizer.SDFARG {
		code = types.InterCode{Content: content, Kind: types.ROW}
	} else if tokens[index].Kind == tokenizer.SASSIGNVARIABLE {
		code = types.InterCode{Content: content, Kind: types.VAR}
	}
	
	functionInterCodeMap[functionPointer] = append(functionInterCodeMap[functionPointer], code)
	index++

	return index, nil
}

// 関数呼び出し文
func functionCall(tokens []tokenizer.Token, index int) (int, error) {
	var err error
	functionCallName := tokens[index].Content

	// 関数名
	if _, ok := functionInterCodeMap[functionCallName]; !ok {
		for name := range functionInterCodeMap {
			fmt.Println(name)
		}
		return index, errors.New(fmt.Sprintf("semantic error: function %s is not defined in line %d", tokens[index].Content, tokens[index].Line))
	}

	index, err = functionName(tokens, index)
	if err != nil {
		return index, err
	}

	// "("
	index++

	// 文字列の並び
	argValues, index, err := rowOfStrings(tokens, index)
	if err != nil {
		return index, err
	}

	functionInterCodeMap[functionPointer] = append(functionInterCodeMap[functionPointer], types.InterCode{Content: functionCallName, Kind: types.CALLFUNC, ArgValues: argValues})

	// ")"
	index++

	return index, nil
}

// 文字列の並び
func rowOfStrings(tokens []tokenizer.Token, index int) ([]string, int, error) {
	var argValues []string
	functionCallName := tokens[index - 2].Content

	// 文字列
	for i, _ := range functionArgMap[functionCallName] {
		if tokens[index].Kind != tokenizer.SSTRING {
			return argValues, index, errors.New(fmt.Sprintf("semantic error: not enough arguments in line %d", tokens[index].Line))
		}

		// functionArgMap[functionCallName][i].Value = tokens[index].Content
		argValues = append(argValues, tokens[index].Content)
		index++

		if i == len(functionArgMap[functionCallName]) - 1 {
			break
		}

		index++
	}

	if tokens[index].Kind == tokenizer.SCOMMA {
		return argValues, index, errors.New(fmt.Sprintf("semantic error: too many arguments in line %d", tokens[index].Line))
	}

	return argValues, index, nil
}

func ifBlock(tokens []tokenizer.Token, index int) (int, error) {
	var err error
	var ifIndexes []int

	// "if"
	index++

	// "("
	index++

	// 条件判定式
	ifContent, index, err := conditionalFormula(tokens, index)
	if err != nil {
		return index, err
	}

	functionInterCodeMap[functionPointer] = append(functionInterCodeMap[functionPointer], types.InterCode{Kind: types.IF, IfContent: ifContent})
	ifIndexes = append(ifIndexes, len(functionInterCodeMap[functionPointer]) - 1)

	// ")"
	index++

	// "{"
	index++

	// 記述部
	index, err = description(tokens, index)
	if err != nil {
		return index, err
	}

	functionInterCodeMap[functionPointer] = append(functionInterCodeMap[functionPointer], types.InterCode{Kind: types.ENDIF})

	// "}"
	index++

	for ;; {
		if tokens[index].Kind != tokenizer.SELIF {
			break
		}

		// elif節
		ifIndexes = append(ifIndexes, len(functionInterCodeMap[functionPointer]))
		index, err = elifSection(tokens, index)
		if err != nil {
			return index, err
		}
	}

	if tokens[index].Kind == tokenizer.SELSE {
		index, err = elseSection(tokens, index)
		if err != nil {
			return index, err
		}
	}

	for _, index := range ifIndexes {
		functionInterCodeMap[functionPointer][index].IfContent.EndIndex = len(functionInterCodeMap[functionPointer])
	}

	return index, nil
}

// 条件判定式
func conditionalFormula(tokens []tokenizer.Token, index int) (types.IfContent, int, error) {
	var ifContent types.IfContent
	var err error
	// 式
	lFormula, index, err := formula(tokens, index)
	if err != nil {
		return ifContent, index, err
	}

	// 比較演算子
	op, index, err := conditionalOperator(tokens, index)
	if err != nil {
		return ifContent, index, err
	}

	// 式
	rFormula, index, err := formula(tokens, index)
	if err != nil {
		return ifContent, index, err
	}

	ifContent = types.IfContent{LFormula: lFormula, RFormula: rFormula, Operator: op}

	return ifContent, index, nil
}

// 式
func formula(tokens []tokenizer.Token, index int) (types.Formula, int, error) {
	// 変数, 文字列
	var formula types.Formula

	if tokens[index].Kind == tokenizer.SIDENTIFIER {
		formula = types.Formula{Content: tokens[index].Content, Kind: tokenizer.SIDENTIFIER}
	} else if tokens[index].Kind == tokenizer.SSTRING {
		formula = types.Formula{Content: tokens[index].Content, Kind: tokenizer.SSTRING}
	}

	index++

	return formula, index, nil
}

// 比較演算子
func conditionalOperator(tokens []tokenizer.Token, index int) (types.OperaterKind, int, error) {
	var op types.OperaterKind
	if tokens[index].Kind == tokenizer.SEQUAL {
		op = types.EQUAL
	} else if tokens[index].Kind == tokenizer.SNOTEQUAL {
		op = types.NOTEQUAL
	}
	// 比較演算子
	index++

	return op, index, nil
}

// elif節
func elifSection(tokens []tokenizer.Token, index int) (int, error) {
	var err error

	// "else if"
	index++

	// "("
	index++

	// 条件判定式
	ifContent, index, err := conditionalFormula(tokens, index)
	if err != nil {
		return index, err
	}

	functionInterCodeMap[functionPointer] = append(functionInterCodeMap[functionPointer], types.InterCode{Kind: types.ELIF, IfContent: ifContent})

	// ")"
	index++

	// "{"
	index++

	// 記述部
	index, err = description(tokens, index)
	if err != nil {
		return index, err
	}

	functionInterCodeMap[functionPointer] = append(functionInterCodeMap[functionPointer], types.InterCode{Kind: types.ENDIF})

	// "}"
	index++
	
	return index, nil
}

// else節
func elseSection(tokens []tokenizer.Token, index int) (int, error) {
	var err error

	// "else"
	functionInterCodeMap[functionPointer] = append(functionInterCodeMap[functionPointer], types.InterCode{Kind: types.ELSE})
	
	index++

	// "{"
	index++

	// 記述部
	index, err = description(tokens, index)
	if err != nil {
		return index, err
	}

	functionInterCodeMap[functionPointer] = append(functionInterCodeMap[functionPointer], types.InterCode{Kind: types.ENDIF})

	// "}"
	index++
	
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
