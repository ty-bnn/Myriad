package compiler

import (
	"fmt"
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

	err = program(newTokens, 0)
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
	argument := types.Argument{Name: name, Kind: types.STRING}

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

	functionArgMap[functionPointer] = append(functionArgMap[functionPointer], argument)

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
	} else if tokens[index].Kind == types.SIF {
		// ifブロック
		index, err = ifBlock(tokens, index)
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
	functionInterCodeMap[functionPointer] = append(functionInterCodeMap[functionPointer], types.InterCode{Content: tokens[index].Content, Kind: types.ROW})
	functionInterCodeMap[functionPointer] = append(functionInterCodeMap[functionPointer], types.InterCode{Content: " ", Kind: types.ROW})
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

	var code types.InterCode
	if tokens[index].Kind == types.SDFARG {
		code = types.InterCode{Content: content, Kind: types.ROW}
	} else if tokens[index].Kind == types.SASSIGNVARIABLE {
		code = types.InterCode{Content: content, Kind: types.VAR}
	}
	
	functionInterCodeMap[functionPointer] = append(functionInterCodeMap[functionPointer], code)
	index++

	return index, nil
}

// 関数呼び出し文
func functionCall(tokens []types.Token, index int) (int, error) {
	var err error
	functionCallName := tokens[index].Content

	// 関数名
	if _, ok := functionInterCodeMap[functionCallName]; !ok {
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
func rowOfStrings(tokens []types.Token, index int) ([]string, int, error) {
	var argValues []string
	functionCallName := tokens[index - 2].Content

	// 文字列
	for i, _ := range functionArgMap[functionCallName] {
		if tokens[index].Kind != types.SSTRING {
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

	if tokens[index].Kind == types.SCOMMA {
		return argValues, index, errors.New(fmt.Sprintf("semantic error: too many arguments in line %d", tokens[index].Line))
	}

	return argValues, index, nil
}

func ifBlock(tokens []types.Token, index int) (int, error) {
	var err error

	// "if"
	index++

	// "("
	index++

	// 条件判定式
	ifContent, index, err := conditionalFormula(tokens, index)
	if err != nil {
		return index, err
	}

	ifCode := &types.InterCode{Kind: types.IF, IfContent: ifContent}
	functionInterCodeMap[functionPointer] = append(functionInterCodeMap[functionPointer], *ifCode)

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
		if tokens[index].Kind != types.SELIF {
			break
		}

		// elif節
		index, err = elifSection(tokens, index)
		if err != nil {
			return index, err
		}
	}

	if tokens[index].Kind == types.SELSE {
		index, err = elseSection(tokens, index)
		if err != nil {
			return index, err
		}
	}

	(*ifCode).IfContent.EndIndex = len(functionInterCodeMap[functionPointer])

	return index, nil
}

// 条件判定式
func conditionalFormula(tokens []types.Token, index int) (types.IfContent, int, error) {
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
func formula(tokens []types.Token, index int) (types.Formula, int, error) {
	// 変数, 文字列
	var formula types.Formula

	if tokens[index].Kind == types.SIDENTIFIER {
		if argumentExist(functionPointer, tokens[index].Content) {
			formula = types.Formula{Content: tokens[index].Content, Kind: types.SIDENTIFIER}
		} else {
			return formula, index, errors.New(fmt.Sprintf("semantic error: function %s is not defined in line %d", tokens[index].Content, tokens[index].Line))
		}
	} else if tokens[index].Kind == types.SSTRING {
		formula = types.Formula{Content: tokens[index].Content, Kind: types.SSTRING}
	}

	index++

	return formula, index, nil
}

// 比較演算子
func conditionalOperator(tokens []types.Token, index int) (types.OperaterKind, int, error) {
	var op types.OperaterKind
	if tokens[index].Kind == types.SEQUAL {
		op = types.EQUAL
	} else if tokens[index].Kind == types.SNOTEQUAL {
		op = types.NOTEQUAL
	}
	// 比較演算子
	index++

	return op, index, nil
}

// elif節
func elifSection(tokens []types.Token, index int) (int, error) {
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
func elseSection(tokens []types.Token, index int) (int, error) {
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