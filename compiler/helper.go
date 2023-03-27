package compiler

import (
	"fmt"
	"os"
	"errors"

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
	var argument Variable

	// 変数名
	name := tokens[index].Content
	if functionPointer == "main" {
		argument = Variable{Name: name, Value:os.Args[argIndex + 3], Kind: ARGUMENT}
	} else {
		argument = Variable{Name: name, Kind: ARGUMENT}
	}

	index, err = variableName(tokens, index)
	if err != nil {
		return index, err
	}

	// "[]"
	if tokens[index].Kind == tokenizer.SARRANGE {
		argument.Kind = VARIABLE
		index++
	}

	functionVarMap[functionPointer] = append(functionVarMap[functionPointer], argument)

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
	} else if tokens[index].Kind == tokenizer.SIDENTIFIER && tokens[index + 1].Kind == tokenizer.SLPAREN {
		// 関数呼び出し文
		index, err = functionCall(tokens, index)
		if err != nil {
			return index, err
		}
	} else if tokens[index].Kind == tokenizer.SIDENTIFIER && tokens[index + 1].Kind == tokenizer.SDEFINE {
		// 変数定義文
		index, err = defineVariable(tokens, index)
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
		functionInterCodeMap[functionPointer] = append(functionInterCodeMap[functionPointer], InterCode{Content: tokens[index].Content, Kind: COMMAND})
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

	var code InterCode
	if tokens[index].Kind == tokenizer.SDFARG {
		code = InterCode{Content: content, Kind: ROW}
	} else if tokens[index].Kind == tokenizer.SASSIGNVARIABLE {
		varIndex, isExist := getVariableIndex(functionPointer, tokens[index].Content)
		if isExist {
			code = InterCode{Content: functionVarMap[functionPointer][varIndex].Value, Kind: ROW}
		} else {
			code = InterCode{Content: content, Kind: VAR}
		}
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
		return index, errors.New(fmt.Sprintf("semantic error: function %s is not defined 1 in line %d", tokens[index].Content, tokens[index].Line))
	}

	index, err = functionName(tokens, index)
	if err != nil {
		return index, err
	}

	// "("
	index++

	// 式の並び
	formulas, index, err := rowOfFormulas(tokens, index)
	if err != nil {
		return index, err
	}

	var newCodes []InterCode
	for _, code := range functionInterCodeMap[functionCallName] {
		var newCode InterCode
		if code.Kind == VAR {
			argIndex, isExist := getArgumentIndex(functionCallName, code.Content)
			if isExist && formulas[argIndex].Kind == tokenizer.SSTRING {
				newCode = InterCode{Content: formulas[argIndex].Content, Kind: ROW}
			} else if isExist && formulas[argIndex].Kind == tokenizer.SIDENTIFIER {
				newCode = InterCode{Content: formulas[argIndex].Content, Kind: VAR}
			} else {
				return index, errors.New(fmt.Sprintf("semantic error: variable %s is not defined 2", code.Content))
			}
		} else if code.Kind == IF || code.Kind == ELIF {
			newCode = code

			if code.IfContent.LFormula.Kind == tokenizer.SIDENTIFIER {
				argIndex, isExist := getArgumentIndex(functionCallName, code.IfContent.LFormula.Content)
				if isExist {
					newCode.IfContent.LFormula = formulas[argIndex]
				} else {
					return index, errors.New(fmt.Sprintf("semantic error: variable %s is not defined 3", code.IfContent.LFormula.Content))
				}
			}

			if code.IfContent.RFormula.Kind == tokenizer.SIDENTIFIER {
				argIndex, isExist := getArgumentIndex(functionCallName, code.IfContent.RFormula.Content)
				if isExist {
					newCode.IfContent.RFormula = formulas[argIndex]
				} else {
					return index, errors.New(fmt.Sprintf("semantic error: variable %s is not defined 4", code.IfContent.RFormula.Content))
				}
			}
		} else {
			newCode = code
		}

		newCodes = append(newCodes, newCode)
	}

	functionInterCodeMap[functionPointer] = append(functionInterCodeMap[functionPointer], newCodes...)

	// ")"
	index++

	return index, nil
}

// 式の並び
func rowOfFormulas(tokens []tokenizer.Token, index int) ([]Formula, int, error) {
	var err error
	var fml Formula
	var formulas []Formula
	functionCallName := tokens[index - 2].Content

	var argNum int
	// 定義された引数の個数を数える
	for _, variable := range functionVarMap[functionCallName] {
		if variable.Kind != ARGUMENT {
			break
		}

		argNum++
	}

	for i := 0; i < argNum; i++ {
		// 式
		if tokens[index].Kind != tokenizer.SSTRING && tokens[index].Kind != tokenizer.SIDENTIFIER {
			return formulas, index, errors.New(fmt.Sprintf("semantic error: not enough arguments in line %d", tokens[index].Line))
		}

		fml, index, err = formula(tokens, index)
		if err != nil {
			return formulas, index, err
		}

		formulas = append(formulas, fml)

		if len(formulas) == argNum {
			break
		}

		// ","
		index++
	}

	if tokens[index].Kind == tokenizer.SCOMMA {
		return formulas, index, errors.New(fmt.Sprintf("semantic error: too many arguments in line %d", tokens[index].Line))
	}

	return formulas, index, nil
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

	functionInterCodeMap[functionPointer] = append(functionInterCodeMap[functionPointer], InterCode{Kind: IF, IfContent: ifContent})
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

	functionInterCodeMap[functionPointer] = append(functionInterCodeMap[functionPointer], InterCode{Kind: ENDIF})

	// "}"
	index++

	for ;; {
		if tokens[index].Kind != tokenizer.SELIF {
			break
		}

		// elif節
		/*
		Note: if節の場合はIFのinterCodeを登録してからfunctionInteCodeMap[functionPointer]の長さを登録しているが、
		      elifとelseの場合はinterCodeを登録する前に長さを登録する必要があるため、-1がいらない         
		*/
		ifIndexes = append(ifIndexes, len(functionInterCodeMap[functionPointer]))
		index, err = elifSection(tokens, index)
		if err != nil {
			return index, err
		}
	}

	if tokens[index].Kind == tokenizer.SELSE {
		ifIndexes = append(ifIndexes, len(functionInterCodeMap[functionPointer]))
		index, err = elseSection(tokens, index)
		if err != nil {
			return index, err
		}
	}

	// ifブロックの最後までのオフセットを格納
	for _, ifIndex := range ifIndexes {
		functionInterCodeMap[functionPointer][ifIndex].IfContent.EndOffset = len(functionInterCodeMap[functionPointer]) - 1 - ifIndex
		fmt.Printf("EndOffset: %d\n", functionInterCodeMap[functionPointer][ifIndex].IfContent.EndOffset)
	}

	// 次のif節（elif, else）までのオフセットを格納
	ifIndexes = append(ifIndexes, len(functionInterCodeMap[functionPointer]))
	for i := 0; i < len(ifIndexes) - 1; i++{
		functionInterCodeMap[functionPointer][ifIndexes[i]].IfContent.NextOffset = ifIndexes[i + 1] - ifIndexes[i] - 1
		fmt.Printf("NextOffset: %d\n", functionInterCodeMap[functionPointer][ifIndexes[i]].IfContent.NextOffset)
	}

	return index, nil
}

// 条件判定式
func conditionalFormula(tokens []tokenizer.Token, index int) (IfContent, int, error) {
	var ifContent IfContent
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

	ifContent = IfContent{LFormula: lFormula, RFormula: rFormula, Operator: op}

	return ifContent, index, nil
}

// 式
func formula(tokens []tokenizer.Token, index int) (Formula, int, error) {
	// 変数, 文字列
	var formula Formula

	if tokens[index].Kind == tokenizer.SIDENTIFIER {
		varIndex, isExist := getVariableIndex(functionPointer, tokens[index].Content)
		if isExist {
			formula = Formula{Content: functionVarMap[functionPointer][varIndex].Value, Kind: tokenizer.SSTRING}
		} else {
			formula = Formula{Content: tokens[index].Content, Kind: tokenizer.SIDENTIFIER}
		}
	} else if tokens[index].Kind == tokenizer.SSTRING {
		formula = Formula{Content: tokens[index].Content, Kind: tokenizer.SSTRING}
	}

	index++

	return formula, index, nil
}

// 比較演算子
func conditionalOperator(tokens []tokenizer.Token, index int) (OperaterKind, int, error) {
	var op OperaterKind
	if tokens[index].Kind == tokenizer.SEQUAL {
		op = EQUAL
	} else if tokens[index].Kind == tokenizer.SNOTEQUAL {
		op = NOTEQUAL
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

	functionInterCodeMap[functionPointer] = append(functionInterCodeMap[functionPointer], InterCode{Kind: ELIF, IfContent: ifContent})

	// ")"
	index++

	// "{"
	index++

	// 記述部
	index, err = description(tokens, index)
	if err != nil {
		return index, err
	}

	functionInterCodeMap[functionPointer] = append(functionInterCodeMap[functionPointer], InterCode{Kind: ENDIF})

	// "}"
	index++
	
	return index, nil
}

// else節
func elseSection(tokens []tokenizer.Token, index int) (int, error) {
	var err error

	// "else"
	functionInterCodeMap[functionPointer] = append(functionInterCodeMap[functionPointer], InterCode{Kind: ELSE})
	
	index++

	// "{"
	index++

	// 記述部
	index, err = description(tokens, index)
	if err != nil {
		return index, err
	}

	functionInterCodeMap[functionPointer] = append(functionInterCodeMap[functionPointer], InterCode{Kind: ENDIF})

	// "}"
	index++
	
	return index, nil
}

// 変数定義文
func defineVariable(tokens []tokenizer.Token, index int) (int, error) {
	var err error

	// 変数名
	name := tokens[index].Content
	index, err = variableName(tokens, index)
	if err != nil {
		return index, err
	}

	// ":="
	index++

	// 文字列
	value := tokens[index].Content
	index++

	newVariable := Variable{Name: name, Value: value, Kind: VARIABLE}
	varIndex, isExist := getVariableIndex(functionPointer, name)
	if isExist {
		functionVarMap[functionPointer][varIndex] = newVariable
	} else {
		functionVarMap[functionPointer] = append(functionVarMap[functionPointer], newVariable)
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
