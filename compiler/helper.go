package compiler

import (
	"fmt"
	"os"
	"errors"

	"myriad/tokenizer"
	"myriad/parser"
	"myriad/helpers"
)

func (c *Compiler) program(tokens []tokenizer.Token, index int) error {
	var err error
	c.tokens = tokens
	c.index = index

	// { 関数インポート文 }
	for c.tokens[c.index].Kind == tokenizer.SIMPORT {
		err = c.importFunc()
		if err != nil {
			return err
		}
	}

	// { 関数 }
	for index < len(c.tokens) {
		if c.tokens[c.index].Kind != tokenizer.SIDENTIFIER {
			break
		}

		err = c.function()
		if err != nil {
			return err
		}
	}
	
	if index >= len(c.tokens) || c.tokens[c.index].Kind != tokenizer.SMAIN {
		return nil
	}

	// メイン部
	err = c.mainFunction()
	if err != nil {
		return err
	}

	return nil
}

// 関数インポート文
func (c *Compiler) importFunc() error {
	var err error
	// "import"
	c.index++

	// 関数名
	err = c.functionName()
	if err != nil {
		return err
	}

	// "from"
	c.index++

	// ファイル名
	err = c.fileName()
	if err != nil {
		return err
	}

	return nil
}

// ファイル名
func (c *Compiler) fileName() error {
	filePath := c.tokens[c.index].Content

	if c.isCompiled(filePath) {
		c.index++
		return nil
	}

	c.readFiles = append(c.readFiles, filePath)

	lines, err := helpers.ReadLinesFromFile(filePath)
	if err != nil {
		return err
	}

	t := &tokenizer.Tokenizer{}
	err = t.Tokenize(lines)
	if err != nil {
		return err
	}

	newTokens := t.Tokens
	// For debug.
	// for _, token := range newTokens {
	// 	fmt.Printf("%30s\t%10d\n", token.Content, token.Kind)
	// }

	p := &parser.Parser{}
	err = p.Parse(newTokens)
	if err != nil {
		return err
	}

	err = c.program(newTokens, 0)
	if err != nil {
		return err
	}

	c.index++

	return nil
}

// 関数
func (c *Compiler) function() error {
	var err error

	// 関数名
	c.functionPointer = c.tokens[c.index].Content

	err = c.functionName()
	if err != nil {
		return err
	}

	// 引数宣言部
	err = c.argumentDecralation()
	if err != nil {
		return err
	}

	// 関数記述部
	err = c.functionDescription()
	if err != nil {
		return err
	}

	return nil
}

// メイン部
func (c *Compiler) mainFunction() error {
	var err error

	c.functionPointer = "main"

	// "main"
	c.index++

	// 引数宣言部
	err = c.argumentDecralation()
	if err != nil {
		return err
	}

	// 関数記述部
	err = c.functionDescription()
	if err != nil {
		return err
	}

	return nil
}

// 引数宣言部
func (c *Compiler) argumentDecralation() error {
	var err error

	// "("
	c.index++

	// 引数群
	if 	c.tokens[c.index].Kind == tokenizer.SIDENTIFIER {
		err = c.arguments()
		if err != nil {
			return err
		}
	}

	// ")"
	c.index++

	return nil
}

// 引数群
func (c *Compiler) arguments() error {
	var err error
	argIndex := 2

	// 変数
	err = c.variable(argIndex)
	if err != nil {
		return err
	}

	for ;; {
		argIndex++
		// ","
		if c.tokens[c.index].Kind != tokenizer.SCOMMA {
			break
		}

		c.index++
		
		// 変数
		err = c.variable(argIndex)
		if err != nil {
			return err
		}
	}

	return nil
}

// 変数
func (c *Compiler) variable(argIndex int) error {
	var err error
	var argument SingleVariable

	// 変数名
	name := c.tokens[c.index].Content
	if c.functionPointer == "main" {
		argument = SingleVariable{VariableCommonDetail: VariableCommonDetail{Name: name, Kind: ARGUMENT}, Value: os.Args[argIndex + 3]}
	} else {
		argument = SingleVariable{VariableCommonDetail: VariableCommonDetail{Name: name, Kind: ARGUMENT}}
	}

	err = c.variableName()
	if err != nil {
		return err
	}

	// "[]"
	if c.tokens[c.index].Kind == tokenizer.SARRANGE {
		argument.Kind = VARIABLE
		c.index++
	}

	c.FunctionVarMap[c.functionPointer] = append(c.FunctionVarMap[c.functionPointer], argument)

	return nil
}

// 関数記述部
func (c *Compiler) functionDescription() error {
	var err error

	// "{"
	c.index++

	// 記述部
	err = c.description()
	if err != nil {
		return err
	}

	// "}"
	c.index++

	return nil
}

// 記述部
func (c *Compiler) description() error {
	var err error

	// 記述ブロック
	err = c.descriptionBlock()
	if err != nil {
		return err
	}

	for ;; {
		if c.tokens[c.index].Kind != tokenizer.SDFCOMMAND && c.tokens[c.index].Kind != tokenizer.SDFARG && c.tokens[c.index].Kind != tokenizer.SIDENTIFIER && c.tokens[c.index].Kind != tokenizer.SIF {
			break
		}

		err = c.descriptionBlock()
		if err != nil {
			return err
		}
	}

	return nil
}

// 記述ブロック
func (c *Compiler) descriptionBlock() error {
	var err error

	if c.tokens[c.index].Kind == tokenizer.SDFCOMMAND || c.tokens[c.index].Kind == tokenizer.SDFARG {
		// Dfile文
		err = c.dockerFile()
		if err != nil {
			return err
		}
	} else if c.tokens[c.index].Kind == tokenizer.SIDENTIFIER && c.tokens[c.index + 1].Kind == tokenizer.SLPAREN {
		// 関数呼び出し文
		err = c.functionCall()
		if err != nil {
			return err
		}
	} else if c.tokens[c.index].Kind == tokenizer.SIDENTIFIER && c.tokens[c.index + 1].Kind == tokenizer.SDEFINE {
		// 変数定義文
		err = c.defineVariable()
		if err != nil {
			return err
		}
	} else if c.tokens[c.index].Kind == tokenizer.SIF {
		// ifブロック
		err = c.ifBlock()
		if err != nil {
			return err
		}
	}

	return nil
}

// Dfile文
func (c *Compiler) dockerFile() error {
	var err error
	if c.tokens[c.index].Kind == tokenizer.SDFCOMMAND {
		// Df命令
		c.FunctionInterCodeMap[c.functionPointer] = append(c.FunctionInterCodeMap[c.functionPointer], InterCode{Content: c.tokens[c.index].Content, Kind: COMMAND})
		c.index++
	}

	// Df引数
	err = c.dfArgs()
	if err != nil {
		return err
	}

	return nil
}

// Df引数部
func (c *Compiler) dfArgs() error {
	var err error
	err = c.dfArg()
	if err != nil {
		return err
	}

	for ;; {
		if c.tokens[c.index].Kind == tokenizer.SDFARG || c.tokens[c.index].Kind == tokenizer.SLDOUBLEBRA {
			err = c.dfArg()
			if err != nil {
				return err
			}
		} else {
			break
		}
	}

	return nil
}

// Df引数
func (c *Compiler) dfArg() error {
	// 生のDF引数
	if c.tokens[c.index].Kind == tokenizer.SDFARG {
		code := InterCode{Content: c.tokens[c.index].Content, Kind: ROW}
		c.FunctionInterCodeMap[c.functionPointer] = append(c.FunctionInterCodeMap[c.functionPointer], code)
		c.index++
	// 置換式
	} else if c.tokens[c.index].Kind == tokenizer.SLDOUBLEBRA {
		err := c.replaceFormula()
		if err != nil {
			return err
		}
	}

	return nil
}

// 置換式
func (c *Compiler) replaceFormula() error {
	// {{
	c.index++

	// 置換変数
	varIndex, isExist := c.getVariableIndex(c.functionPointer, c.tokens[c.index].Content)

	var code InterCode
	if isExist {
		value, err := c.FunctionVarMap[c.functionPointer][varIndex].getValue(0)
		if err != nil {
			return err
		}

		code = InterCode{Content: value, Kind: ROW}
	} else {
		code = InterCode{Content: c.tokens[c.index].Content, Kind: VAR}
	}

	c.FunctionInterCodeMap[c.functionPointer] = append(c.FunctionInterCodeMap[c.functionPointer], code)
	c.index++

	// }}
	c.index++

	return nil
}

// 関数呼び出し文
func (c *Compiler) functionCall() error {
	var err error
	functionCallName := c.tokens[c.index].Content

	// 関数名
	if _, ok := c.FunctionInterCodeMap[functionCallName]; !ok {
		return errors.New(fmt.Sprintf("semantic error: function %s is not defined 1 in line %d", c.tokens[c.index].Content, c.tokens[c.index].Line))
	}

	err = c.functionName()
	if err != nil {
		return err
	}

	// "("
	c.index++

	// 式の並び
	formulas, err := c.rowOfFormulas()
	if err != nil {
		return err
	}

	var newCodes []InterCode
	for _, code := range c.FunctionInterCodeMap[functionCallName] {
		var newCode InterCode
		if code.Kind == VAR {
			argIndex, isExist := c.getArgumentIndex(functionCallName, code.Content)
			if isExist && formulas[argIndex].Kind == tokenizer.SSTRING {
				newCode = InterCode{Content: formulas[argIndex].Content, Kind: ROW}
			} else if isExist && formulas[argIndex].Kind == tokenizer.SIDENTIFIER {
				newCode = InterCode{Content: formulas[argIndex].Content, Kind: VAR}
			} else {
				return errors.New(fmt.Sprintf("semantic error: variable %s is not defined 2", code.Content))
			}
		} else if code.Kind == IF || code.Kind == ELIF {
			newCode = code

			if code.IfContent.LFormula.Kind == tokenizer.SIDENTIFIER {
				argIndex, isExist := c.getArgumentIndex(functionCallName, code.IfContent.LFormula.Content)
				if isExist {
					newCode.IfContent.LFormula = formulas[argIndex]
				} else {
					return errors.New(fmt.Sprintf("semantic error: variable %s is not defined 3", code.IfContent.LFormula.Content))
				}
			}

			if code.IfContent.RFormula.Kind == tokenizer.SIDENTIFIER {
				argIndex, isExist := c.getArgumentIndex(functionCallName, code.IfContent.RFormula.Content)
				if isExist {
					newCode.IfContent.RFormula = formulas[argIndex]
				} else {
					return errors.New(fmt.Sprintf("semantic error: variable %s is not defined 4", code.IfContent.RFormula.Content))
				}
			}
		} else {
			newCode = code
		}

		newCodes = append(newCodes, newCode)
	}

	c.FunctionInterCodeMap[c.functionPointer] = append(c.FunctionInterCodeMap[c.functionPointer], newCodes...)

	// ")"
	c.index++

	return nil
}

// 式の並び
func (c *Compiler) rowOfFormulas() ([]Formula, error) {
	var err error
	var fml Formula
	var formulas []Formula
	functionCallName := c.tokens[c.index - 2].Content

	var argNum int
	// 定義された引数の個数を数える
	for _, variable := range c.FunctionVarMap[functionCallName] {
		if variable.getKind() != ARGUMENT {
			break
		}

		argNum++
	}

	for i := 0; i < argNum; i++ {
		// 式
		if c.tokens[c.index].Kind != tokenizer.SSTRING && c.tokens[c.index].Kind != tokenizer.SIDENTIFIER {
			return formulas, errors.New(fmt.Sprintf("semantic error: not enough arguments in line %d", c.tokens[c.index].Line))
		}

		fml, err = c.formula()
		if err != nil {
			return formulas, err
		}

		formulas = append(formulas, fml)

		if len(formulas) == argNum {
			break
		}

		// ","
		c.index++
	}

	if c.tokens[c.index].Kind == tokenizer.SCOMMA {
		return formulas, errors.New(fmt.Sprintf("semantic error: too many arguments in line %d", c.tokens[c.index].Line))
	}

	return formulas, nil
}

func (c *Compiler) ifBlock() error {
	var err error
	var ifIndexes []int

	// "if"
	c.index++

	// "("
	c.index++

	// 条件判定式
	ifContent, err := c.conditionalFormula()
	if err != nil {
		return err
	}

	c.FunctionInterCodeMap[c.functionPointer] = append(c.FunctionInterCodeMap[c.functionPointer], InterCode{Kind: IF, IfContent: ifContent})
	ifIndexes = append(ifIndexes, len(c.FunctionInterCodeMap[c.functionPointer]) - 1)

	// ")"
	c.index++

	// "{"
	c.index++

	// 記述部
	err = c.description()
	if err != nil {
		return err
	}

	c.FunctionInterCodeMap[c.functionPointer] = append(c.FunctionInterCodeMap[c.functionPointer], InterCode{Kind: ENDIF})

	// "}"
	c.index++

	for ;; {
		if c.tokens[c.index].Kind != tokenizer.SELIF {
			break
		}

		// elif節
		/*
		Note: if節の場合はIFのinterCodeを登録してからfunctionInteCodeMap[c.functionPointer]の長さを登録しているが、
		      elifとelseの場合はinterCodeを登録する前に長さを登録する必要があるため、-1がいらない         
		*/
		ifIndexes = append(ifIndexes, len(c.FunctionInterCodeMap[c.functionPointer]))
		err = c.elifSection()
		if err != nil {
			return err
		}
	}

	if c.tokens[c.index].Kind == tokenizer.SELSE {
		ifIndexes = append(ifIndexes, len(c.FunctionInterCodeMap[c.functionPointer]))
		err = c.elseSection()
		if err != nil {
			return err
		}
	}

	// ifブロックの最後までのオフセットを格納
	for _, ifIndex := range ifIndexes {
		c.FunctionInterCodeMap[c.functionPointer][ifIndex].IfContent.EndOffset = len(c.FunctionInterCodeMap[c.functionPointer]) - 1 - ifIndex
		fmt.Printf("EndOffset: %d\n", c.FunctionInterCodeMap[c.functionPointer][ifIndex].IfContent.EndOffset)
	}

	// 次のif節（elif, else）までのオフセットを格納
	ifIndexes = append(ifIndexes, len(c.FunctionInterCodeMap[c.functionPointer]))
	for i := 0; i < len(ifIndexes) - 1; i++{
		c.FunctionInterCodeMap[c.functionPointer][ifIndexes[i]].IfContent.NextOffset = ifIndexes[i + 1] - ifIndexes[i] - 1
		fmt.Printf("NextOffset: %d\n", c.FunctionInterCodeMap[c.functionPointer][ifIndexes[i]].IfContent.NextOffset)
	}

	return nil
}

// 条件判定式
func (c *Compiler) conditionalFormula() (IfContent, error) {
	var ifContent IfContent
	var err error
	// 式
	lFormula, err := c.formula()
	if err != nil {
		return ifContent, err
	}

	// 比較演算子
	op, err := c.conditionalOperator()
	if err != nil {
		return ifContent, err
	}

	// 式
	rFormula, err := c.formula()
	if err != nil {
		return ifContent, err
	}

	ifContent = IfContent{LFormula: lFormula, RFormula: rFormula, Operator: op}

	return ifContent, nil
}

// 式
func (c *Compiler) formula() (Formula, error) {
	// 変数, 文字列
	var formula Formula

	if c.tokens[c.index].Kind == tokenizer.SIDENTIFIER {
		varIndex, isExist := c.getVariableIndex(c.functionPointer, c.tokens[c.index].Content)
		if isExist {
			value, err := c.FunctionVarMap[c.functionPointer][varIndex].getValue(0)
			if err != nil {
				return Formula{}, err
			}
			formula = Formula{Content: value, Kind: tokenizer.SSTRING}
		} else {
			formula = Formula{Content: c.tokens[c.index].Content, Kind: tokenizer.SIDENTIFIER}
		}
	} else if c.tokens[c.index].Kind == tokenizer.SSTRING {
		formula = Formula{Content: c.tokens[c.index].Content, Kind: tokenizer.SSTRING}
	}

	c.index++

	return formula, nil
}

// 比較演算子
func (c *Compiler) conditionalOperator() (OperaterKind, error) {
	var op OperaterKind
	if c.tokens[c.index].Kind == tokenizer.SEQUAL {
		op = EQUAL
	} else if c.tokens[c.index].Kind == tokenizer.SNOTEQUAL {
		op = NOTEQUAL
	}
	// 比較演算子
	c.index++

	return op, nil
}

// elif節
func (c *Compiler) elifSection() error {
	var err error

	// "else if"
	c.index++

	// "("
	c.index++

	// 条件判定式
	ifContent, err := c.conditionalFormula()
	if err != nil {
		return err
	}

	c.FunctionInterCodeMap[c.functionPointer] = append(c.FunctionInterCodeMap[c.functionPointer], InterCode{Kind: ELIF, IfContent: ifContent})

	// ")"
	c.index++

	// "{"
	c.index++

	// 記述部
	err = c.description()
	if err != nil {
		return err
	}

	c.FunctionInterCodeMap[c.functionPointer] = append(c.FunctionInterCodeMap[c.functionPointer], InterCode{Kind: ENDIF})

	// "}"
	c.index++
	
	return nil
}

// else節
func (c *Compiler) elseSection() error {
	var err error

	// "else"
	c.FunctionInterCodeMap[c.functionPointer] = append(c.FunctionInterCodeMap[c.functionPointer], InterCode{Kind: ELSE})
	
	c.index++

	// "{"
	c.index++

	// 記述部
	err = c.description()
	if err != nil {
		return err
	}

	c.FunctionInterCodeMap[c.functionPointer] = append(c.FunctionInterCodeMap[c.functionPointer], InterCode{Kind: ENDIF})

	// "}"
	c.index++
	
	return nil
}

// 変数定義文
func (c *Compiler) defineVariable() error {
	var err error

	// 変数名
	name := c.tokens[c.index].Content
	err = c.variableName()
	if err != nil {
		return err
	}

	// ":="
	c.index++

	var newVariable Variable
	if c.tokens[c.index].Kind == tokenizer.SSTRING {
		// 文字列
		value := c.tokens[c.index].Content
		c.index++
		newVariable = SingleVariable{VariableCommonDetail: VariableCommonDetail{Name: name, Kind: VARIABLE}, Value: value}
	} else if c.tokens[c.index].Kind == tokenizer.SLBRACE {
		// 配列
		values, err := c.array()
		if err != nil {
			return err
		}
		newVariable = MultiVariable{VariableCommonDetail: VariableCommonDetail{Name: name, Kind: VARIABLE}, Values: values}
	}

	varIndex, isExist := c.getVariableIndex(c.functionPointer, name)
	if isExist {
		// 変数を書き換える, 後に定義された値を使用する
		c.FunctionVarMap[c.functionPointer][varIndex] = newVariable
	} else {
		// 新しく変数を登録する
		c.FunctionVarMap[c.functionPointer] = append(c.FunctionVarMap[c.functionPointer], newVariable)
	}

	return nil
}

// 配列
func (c *Compiler) array() ([]string, error) {
	// {
	c.index++
	
	// 文字列の並び
	values, err := c.rowOfStrings()
	if err != nil {
		return nil, err
	}

	// }
	c.index++

	return values, nil
}

// 文字列の並び
func (c *Compiler) rowOfStrings() ([]string, error) {
	var strings []string
	// 文字列
	strings = append(strings, c.tokens[c.index].Content)
	c.index++

	for ;; {
		// ","
		if c.tokens[c.index].Kind != tokenizer.SCOMMA {
			break
		}

		c.index++


		// 文字列
		strings = append(strings, c.tokens[c.index].Content)
		c.index++
	}

	return strings, nil
}

// 関数名
func (c *Compiler) functionName() error {
	// 名前
	c.index++

	return nil
}

// 変数名
func (c *Compiler) variableName() error {
	// 名前
	c.index++

	return nil
}
