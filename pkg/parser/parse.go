package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/ty-bnn/myriad/pkg/model/codes"
	"github.com/ty-bnn/myriad/pkg/model/token"
	"github.com/ty-bnn/myriad/pkg/model/values"
	"github.com/ty-bnn/myriad/pkg/tokenizer"
	"github.com/ty-bnn/myriad/pkg/utils"
)

func (p *Parser) Parse() error {
	fmt.Println("Parsing...")

	err := p.program()
	if err != nil {
		return err
	}

	fmt.Println("Parse Done.")
	return nil
}

func (p *Parser) program() error {
	var err error
	// { 関数インポート文 }
	for {
		if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.IMPORT {
			break
		}

		err = p.importFunc()
		if err != nil {
			return err
		}
	}

	// { 関数 }
	for {
		if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.IDENTIFIER {
			break
		}

		err = p.function()
		if err != nil {
			return err
		}
	}

	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.MAIN {
		return nil
	}

	// メイン部
	err = p.mainFunction()
	if err != nil {
		return err
	}

	return nil
}

// 関数インポート文
func (p *Parser) importFunc() error {
	var err error
	// "import"
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.IMPORT {
		return errors.New(fmt.Sprintf("syntax error: cannot find 'import'"))
	}

	p.index++

	// 関数名
	_, err = p.functionName()
	if err != nil {
		return err
	}

	// "from"
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.FROM {
		return errors.New(fmt.Sprintf("syntax error: cannot find 'from'"))
	}

	p.index++

	// ファイル名
	filePath, err := p.fileName()
	if err != nil {
		return err
	}

	if p.isCompiled(filePath) {
		return nil
	}

	lines, err := utils.ReadLinesFromFile(filePath)
	if err != nil {
		return err
	}

	t := tokenizer.New(lines)
	err = t.Tokenize()
	if err != nil {
		return err
	}

	newP := New(t.Tokens)
	err = newP.Parse()
	if err != nil {
		return err
	}

	err = p.addFuncCodes(newP.FuncToCodes)
	if err != nil {
		return err
	}

	p.compiledFiles = append(p.compiledFiles, filePath)

	return nil
}

// ファイル名
func (p *Parser) fileName() (string, error) {
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.STRING {
		return "", errors.New(fmt.Sprintf("syntax error: cannot find an identifier"))
	}

	fileName := p.tokens[p.index].Content

	p.index++

	return fileName, nil
}

// TODO: 関数もmainもやる処理は変わらないのでくっつける
// 関数
func (p *Parser) function() error {
	var err error

	// 関数名
	funcName, err := p.functionName()
	if err != nil {
		return err
	}

	if _, has := p.FuncToCodes[funcName]; has {
		return errors.New(fmt.Sprintf("semantic error: %s is already declared", funcName))
	}

	// 引数宣言部
	argCodes, err := p.argumentDeclaration()
	if err != nil {
		return err
	}

	// 関数記述部
	funcCodes, err := p.functionDescription()
	if err != nil {
		return err
	}

	p.FuncToCodes[funcName] = append(p.FuncToCodes[funcName], argCodes...)
	p.FuncToCodes[funcName] = append(p.FuncToCodes[funcName], funcCodes...)

	return nil
}

// メイン部
func (p *Parser) mainFunction() error {
	var err error

	// "main"
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.MAIN {
		return errors.New(fmt.Sprintf("syntax error: cannot find 'main'"))
	}

	funcName := "main"
	if _, has := p.FuncToCodes[funcName]; has {
		return errors.New(fmt.Sprintf("semantic error: %s is already declared", funcName))
	}

	p.index++

	// 引数宣言部
	argCodes, err := p.argumentDeclaration()
	if err != nil {
		return err
	}

	// 関数記述部
	funcCodes, err := p.functionDescription()
	if err != nil {
		return err
	}

	p.FuncToCodes[funcName] = append(p.FuncToCodes[funcName], argCodes...)
	p.FuncToCodes[funcName] = append(p.FuncToCodes[funcName], funcCodes...)

	return nil
}

// 引数宣言部
func (p *Parser) argumentDeclaration() ([]codes.Code, error) {
	var argCodes []codes.Code
	var err error

	// "("
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.LPAREN {
		return nil, errors.New(fmt.Sprintf("syntax error: cannot find '('"))
	}

	p.index++

	// 引数群
	if p.index < len(p.tokens) && p.tokens[p.index].Kind == token.IDENTIFIER {
		argCodes, err = p.arguments()
		if err != nil {
			return nil, err
		}
	}

	// ")"
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.RPAREN {
		return nil, errors.New(fmt.Sprintf("syntax error: cannot find ')'"))
	}

	p.index++

	return argCodes, nil
}

// 引数群
func (p *Parser) arguments() ([]codes.Code, error) {
	var argCodes []codes.Code
	var err error

	// 変数名
	argName, err := p.variableName()
	if err != nil {
		return nil, err
	}

	// TODO: Valueをnilのままにしたい
	defCode := codes.Define{
		Kind: codes.DEFINE,
		Key:  argName,
	}

	argCodes = append(argCodes, defCode)

	for {
		// ","
		if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.COMMA {
			break
		}

		p.index++

		// 変数名
		argName, err = p.variableName()
		if err != nil {
			return nil, err
		}

		// TODO: Valueをnilのままにしたい
		defCode = codes.Define{
			Kind: codes.DEFINE,
			Key:  argName,
		}

		argCodes = append(argCodes, defCode)
	}

	return argCodes, nil
}

// 関数記述部
// TODO: if文やfor文の処理記述部分にも適用可能
func (p *Parser) functionDescription() ([]codes.Code, error) {
	var err error

	// "{"
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.LBRACE {
		return nil, errors.New(fmt.Sprintf("syntax error: cannot find '{'"))
	}

	p.index++

	// 記述部
	descCodes, err := p.description()
	if err != nil {
		return nil, err
	}

	// "}"
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.RBRACE {
		return nil, errors.New(fmt.Sprintf("syntax error: cannot find '}'"))
	}

	p.index++

	return descCodes, nil
}

// 記述部
func (p *Parser) description() ([]codes.Code, error) {
	var descCodes []codes.Code
	var err error

	// 記述ブロック
	descBCodes, err := p.descriptionBlock()
	if err != nil {
		return nil, err
	}

	descCodes = append(descCodes, descBCodes...)

	for {
		if p.index >= len(p.tokens) || (p.tokens[p.index].Kind != token.DFBEGIN && p.tokens[p.index].Kind != token.DFARG && p.tokens[p.index].Kind != token.IDENTIFIER && p.tokens[p.index].Kind != token.IF && p.tokens[p.index].Kind != token.FOR) {
			break
		}

		descBCodes, err := p.descriptionBlock()
		if err != nil {
			return nil, err
		}

		descCodes = append(descCodes, descBCodes...)
	}

	return descCodes, nil
}

// 記述ブロック
func (p *Parser) descriptionBlock() ([]codes.Code, error) {
	if p.index < len(p.tokens) && p.tokens[p.index].Kind == token.DFBEGIN {
		// Dfileブロック
		dockerCodes, err := p.dockerfileBlock()
		if err != nil {
			return nil, err
		}

		return dockerCodes, nil
	} else if p.index+1 < len(p.tokens) && p.tokens[p.index].Kind == token.IDENTIFIER && p.tokens[p.index+1].Kind == token.LPAREN {
		// 関数呼び出し文
		funcCode, err := p.functionCall()
		if err != nil {
			return nil, err
		}

		return []codes.Code{funcCode}, err
	} else if p.index+1 < len(p.tokens) && p.tokens[p.index].Kind == token.IDENTIFIER && p.tokens[p.index+1].Kind == token.DEFINE {
		// 変数定義文
		dvCode, err := p.defineVariable()
		if err != nil {
			return nil, err
		}

		return []codes.Code{dvCode}, err
	} else if p.index+1 < len(p.tokens) && p.tokens[p.index].Kind == token.IDENTIFIER && p.tokens[p.index+1].Kind == token.ASSIGN {
		// 変数代入文
		avCode, err := p.assignVariable()
		if err != nil {
			return nil, err
		}

		return []codes.Code{avCode}, err
	} else if p.index < len(p.tokens) && p.tokens[p.index].Kind == token.IF {
		// ifブロック
		ifCodes, err := p.ifBlock()
		if err != nil {
			return nil, err
		}

		return ifCodes, nil
	} else if p.index < len(p.tokens) && p.tokens[p.index].Kind == token.FOR {
		// forブロック
		forCodes, err := p.forBlock()
		if err != nil {
			return nil, err
		}

		return forCodes, nil
	}

	return nil, errors.New(fmt.Sprintf("syntax error: cannot find a description block"))
}

// Dfileブロック
func (p *Parser) dockerfileBlock() ([]codes.Code, error) {
	var dockerBlockCodes []codes.Code

	// {{-
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.DFBEGIN {
		return nil, errors.New(fmt.Sprintf("syntax error: cannot find '{{-'"))
	}
	p.index++

	for {
		if p.index >= len(p.tokens) || (p.tokens[p.index].Kind != token.DFCOMMAND && p.tokens[p.index].Kind != token.DFARG && p.tokens[p.index].Kind != token.LDOUBLEBRA) {
			break
		}

		dockerCodes, err := p.dockerFile()
		if err != nil {
			return nil, err
		}
		dockerBlockCodes = append(dockerBlockCodes, dockerCodes...)
	}

	// -}}
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.DFEND {
		return nil, errors.New(fmt.Sprintf("syntax error: cannot find '-}}'"))
	}
	p.index++

	return dockerBlockCodes, nil
}

// Dfile文
func (p *Parser) dockerFile() ([]codes.Code, error) {
	var dockerCodes []codes.Code
	var err error

	// Df命令
	if p.tokens[p.index].Kind == token.DFCOMMAND {
		cmdCode := codes.Command{Kind: codes.COMMAND, Content: p.tokens[p.index].Content}
		dockerCodes = append(dockerCodes, cmdCode)
		p.index++
	}

	// Df引数部
	dfArgsCodes, err := p.dfArgs()
	if err != nil {
		return nil, err
	}

	dockerCodes = append(dockerCodes, dfArgsCodes...)

	return dockerCodes, nil
}

// Df引数部
func (p *Parser) dfArgs() ([]codes.Code, error) {
	var dfArgsCodes []codes.Code
	var err error

	dfArgCode, err := p.dfArg()
	if err != nil {
		return nil, err
	}

	dfArgsCodes = append(dfArgsCodes, dfArgCode)

	for {
		if p.index < len(p.tokens) && (p.tokens[p.index].Kind == token.DFARG || p.tokens[p.index].Kind == token.LDOUBLEBRA) {
			dfArgCode, err := p.dfArg()
			if err != nil {
				return nil, err
			}

			dfArgsCodes = append(dfArgsCodes, dfArgCode)
		} else {
			break
		}
	}

	return dfArgsCodes, nil
}

// Df引数
func (p *Parser) dfArg() (codes.Code, error) {
	if p.index < len(p.tokens) && p.tokens[p.index].Kind == token.DFARG {
		// 生のDf引数
		rowCode := codes.Literal{Kind: codes.LITERAL, Content: p.tokens[p.index].Content}
		p.index++

		return rowCode, nil
	} else if p.index < len(p.tokens) && p.tokens[p.index].Kind == token.LDOUBLEBRA {
		//置換式
		repCode, err := p.replaceFormula()
		if err != nil {
			return nil, err
		}

		return repCode, nil
	}

	return nil, errors.New(fmt.Sprintf("syntax error: cannot find Df argument"))
}

// 置換式
func (p *Parser) replaceFormula() (codes.Code, error) {
	// {{
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.LDOUBLEBRA {
		return nil, errors.New(fmt.Sprintf("syntax error: cannot find '{{'"))
	}
	p.index++

	target, err := p.singleAssignFormula()
	if err != nil {
		return nil, err
	}
	repCode := codes.Replace{Kind: codes.REPLACE, Value: target}

	// }}
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.RDOUBLEBRA {
		return nil, errors.New(fmt.Sprintf("syntax error: cannot find '}}'"))
	}

	p.index++

	return repCode, nil
}

// 関数呼び出し文
func (p *Parser) functionCall() (codes.Code, error) {
	var args []values.Value
	var err error

	// 関数名
	funcName, err := p.functionName()
	if err != nil {
		return nil, err
	}

	// "("
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.LPAREN {
		return nil, errors.New(fmt.Sprintf("syntax error: cannot find '('"))
	}

	p.index++

	// 代入値の並び
	if p.index < len(p.tokens) && (p.tokens[p.index].Kind == token.STRING || p.tokens[p.index].Kind == token.IDENTIFIER ||
		p.tokens[p.index].Kind == token.LBRACE || p.tokens[p.index].Kind == token.JSONUNMARSHAL) {
		args, err = p.rowOfAssignValues()
		if err != nil {
			return nil, err
		}
	}

	// ")"
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.RPAREN {
		return nil, errors.New(fmt.Sprintf("syntax error: cannot find ')'"))
	}

	cpCode := codes.CallProc{Kind: codes.CALLPROC, ProcName: funcName, Args: args}

	p.index++

	return cpCode, nil
}

// 代入値の並び
func (p *Parser) rowOfAssignValues() ([]values.Value, error) {
	var values []values.Value
	var err error
	// 代入値
	value, err := p.assignValue()
	if err != nil {
		return nil, err
	}
	values = append(values, value)

	for {
		// ","
		if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.COMMA {
			break
		}
		p.index++
		// 代入値
		value, err := p.assignValue()
		if err != nil {
			return values, err
		}
		values = append(values, value)
	}

	return values, nil
}

// ifブロック
func (p *Parser) ifBlock() ([]codes.Code, error) {
	var ifBCodes []codes.Code
	var err error

	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.IF {
		return nil, errors.New(fmt.Sprintf("syntax error: cannot find 'if'"))
	}

	ifCodes, err := p.ifSection()
	if err != nil {
		return nil, err
	}

	ifBCodes = append(ifBCodes, ifCodes...)

	for {
		if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.ELIF {
			break
		}

		// elif節
		elifCodes, err := p.elifSection()
		if err != nil {
			return nil, err
		}

		ifBCodes = append(ifBCodes, elifCodes...)
	}

	if p.index < len(p.tokens) && p.tokens[p.index].Kind == token.ELSE {
		elseCodes, err := p.elseSection()
		if err != nil {
			return nil, err
		}
		ifBCodes = append(ifBCodes, elseCodes...)
	}

	return ifBCodes, nil
}

// 条件判定式
func (p *Parser) conditionalFormula() (*codes.ConditionalNode, error) {
	// 項
	root, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.index < len(p.tokens) && p.tokens[p.index].Kind == token.OR {
		// ||
		p.index++

		// 項
		rNode, err := p.term()
		if err != nil {
			return nil, err
		}

		newRoot := codes.ConditionalNode{
			Operator: codes.OR,
			Left:     root,
			Right:    rNode,
		}

		root = &newRoot
	}

	return root, nil
}

// 項
func (p *Parser) term() (*codes.ConditionalNode, error) {
	// 因子
	root, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.index < len(p.tokens) && p.tokens[p.index].Kind == token.AND {
		// &&
		p.index++

		// 因子
		rNode, err := p.factor()
		if err != nil {
			return nil, err
		}

		newRoot := codes.ConditionalNode{
			Operator: codes.AND,
			Left:     root,
			Right:    rNode,
		}

		root = &newRoot
	}

	return root, nil
}

// 因子
func (p *Parser) factor() (*codes.ConditionalNode, error) {
	if p.index < len(p.tokens) && p.tokens[p.index].Kind == token.LPAREN {
		// (
		p.index++

		// 条件判定式
		condFml, err := p.conditionalFormula()
		if err != nil {
			return nil, err
		}

		// )
		if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.RPAREN {
			return nil, errors.New(fmt.Sprintf("syntax error: cannot find )"))
		}
		p.index++

		return condFml, nil
	}

	compFml, err := p.compFormula()
	if err != nil {
		return nil, err
	}

	return compFml, nil
}

// 比較式
func (p *Parser) compFormula() (*codes.ConditionalNode, error) {
	var err error
	left, err := p.singleAssignFormula()
	if err != nil {
		return nil, err
	}
	lNode := codes.ConditionalNode{Var: left}

	op, err := p.conditionalOperator()
	if err != nil {
		return nil, err
	}

	right, err := p.singleAssignFormula()
	if err != nil {
		return nil, err
	}
	rNode := codes.ConditionalNode{Var: right}

	return &codes.ConditionalNode{Operator: op, Left: &lNode, Right: &rNode}, nil
}

// 比較演算子
func (p *Parser) conditionalOperator() (codes.OperatorKind, error) {
	var op codes.OperatorKind

	// "==", "!="
	if p.index >= len(p.tokens) || (p.tokens[p.index].Kind != token.EQUAL && p.tokens[p.index].Kind != token.NOTEQUAL) {
		return -1, errors.New(fmt.Sprintf("syntax error: cannot find conditional operator"))
	}

	if p.tokens[p.index].Kind == token.EQUAL {
		op = codes.EQUAL
	} else if p.tokens[p.index].Kind == token.NOTEQUAL {
		op = codes.NOTEQUAL
	}

	p.index++

	return op, nil
}

// 変数定義文
func (p *Parser) defineVariable() (codes.Code, error) {
	var err error

	// 変数名
	vName, err := p.variableName()
	if err != nil {
		return nil, err
	}

	// ":="
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.DEFINE {
		return nil, errors.New(fmt.Sprintf("syntax error: cannot find ':='"))
	}

	p.index++

	value, err := p.assignValue()
	if err != nil {
		return nil, err
	}

	defCode := codes.Define{Kind: codes.DEFINE, Key: vName, Value: value}

	return defCode, nil
}

// 変数代入文
func (p *Parser) assignVariable() (codes.Code, error) {
	// 変数名
	vName, err := p.variableName()
	if err != nil {
		return nil, err
	}

	// "="
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.ASSIGN {
		return nil, errors.New(fmt.Sprintf("syntax error: cannot find '='"))
	}

	p.index++

	value, err := p.assignValue()
	if err != nil {
		return nil, err
	}

	defCode := codes.Assign{Kind: codes.ASSIGN, Key: vName, Value: value}

	return defCode, nil
}

// 代入値
func (p *Parser) assignValue() (values.Value, error) {
	if p.tokenIs(token.LBRACE, 0) || p.tokenIs(token.JSONUNMARSHAL, 0) || (p.tokenIs(token.IDENTIFIER, 0) && p.tokenIs(token.DOT, 1)) {
		value, err := p.complexAssignValue()
		return value, err
	}

	value, err := p.singleAssignFormula()
	return value, err
}

// 単一代入式
func (p *Parser) singleAssignFormula() (values.Value, error) {
	value, err := p.singleAssignValue()
	if err != nil {
		return nil, err
	}

	if !p.tokenIs(token.PLUS, 0) {
		return value, nil
	}

	var vls []values.Value
	vls = append(vls, value)

	for p.tokenIs(token.PLUS, 0) {
		// +
		p.index++

		value, err = p.singleAssignValue()
		if err != nil {
			return nil, err
		}
		vls = append(vls, value)
	}

	return values.AddString{Kind: values.ADDSTRING, Values: vls}, nil
}

// 単一代入値
func (p *Parser) singleAssignValue() (values.Value, error) {
	if p.tokenIs(token.STRING, 0) {
		value := values.Literal{Kind: values.LITERAL, Value: p.tokens[p.index].Content}
		p.index++
		return value, nil
	} else if p.tokenIs(token.IDENTIFIER, 0) && p.tokenIs(token.LBRACKET, 1) && p.tokenIs(token.NUMBER, 2) {
		value, err := p.arrayElement()
		return value, err
	} else if p.tokenIs(token.IDENTIFIER, 0) && p.tokenIs(token.LBRACKET, 1) {
		value, err := p.mapValue()
		return value, err
	} else if p.tokenIs(token.IDENTIFIER, 0) {
		vName, err := p.variableName()
		value := values.Ident{Kind: values.IDENT, Name: vName}
		return value, err
	}
	return nil, errors.New(fmt.Sprintf("syntax error: cannot find complex assign value"))
}

// 複合代入値
func (p *Parser) complexAssignValue() (values.Value, error) {
	if p.tokenIs(token.LBRACE, 0) {
		// 配列
		arrValues, err := p.array()
		return values.Literals{Kind: values.LITERALS, Values: arrValues}, err
	} else if p.tokenIs(token.JSONUNMARSHAL, 0) {
		// JsonUnmarshal
		jsonData, err := p.jsonUnmarshal()
		return values.Map{Kind: values.MAP, Value: jsonData}, err
	} else if p.tokenIs(token.IDENTIFIER, 0) {
		// mapキー
		value, err := p.mapKey()
		return value, err
	}
	return nil, errors.New(fmt.Sprintf("syntax error: cannot find complex assign value"))
}

// 配列
func (p *Parser) array() ([]string, error) {
	// {
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.LBRACE {
		return []string{}, errors.New(fmt.Sprintf("syntax error: cannot find '{'"))
	}

	p.index++

	arrayValues, err := p.rowOfStrings()
	if err != nil {
		return arrayValues, err
	}

	// }
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.RBRACE {
		return nil, errors.New(fmt.Sprintf("syntax error: cannot find '}'"))
	}

	p.index++

	return arrayValues, nil
}

// 配列要素
func (p *Parser) arrayElement() (values.Element, error) {
	// 変数名
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.IDENTIFIER {
		return values.Element{}, errors.New(fmt.Sprintf("syntax error: cannot find identifier"))
	}

	name := p.tokens[p.index].Content
	p.index++

	// [
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.LBRACKET {
		return values.Element{}, errors.New(fmt.Sprintf("syntax error: cannot find left bracket"))
	}

	p.index++

	// 数字
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.NUMBER {
		return values.Element{}, errors.New(fmt.Sprintf("syntax error: cannot find number"))
	}

	index, err := strconv.Atoi(p.tokens[p.index].Content)
	if err != nil {
		return values.Element{}, err
	}
	p.index++

	// [
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.RBRACKET {
		return values.Element{}, errors.New(fmt.Sprintf("syntax error: cannot find right bracket"))
	}

	p.index++

	return values.Element{Kind: values.ELEMENT, Name: name, Index: index}, nil
}

// 文字列の並び
func (p *Parser) rowOfStrings() ([]string, error) {
	var strings []string

	// 文字列
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.STRING {
		return nil, errors.New(fmt.Sprintf("syntax error: cannot find string"))
	}

	strings = append(strings, p.tokens[p.index].Content)

	p.index++

	for {
		// ","
		if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.COMMA {
			break
		}

		p.index++

		// 文字列
		if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.STRING {
			return strings, errors.New(fmt.Sprintf("syntax error: cannot find string"))
		}

		strings = append(strings, p.tokens[p.index].Content)

		p.index++
	}

	return strings, nil
}

// Json読み取り
func (p *Parser) jsonUnmarshal() (map[string]interface{}, error) {
	// JsonUnmarshal
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.JSONUNMARSHAL {
		return nil, errors.New(fmt.Sprintf("syntax error: cannot find JsonUnmarshal"))
	}
	p.index++

	// (
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.LPAREN {
		return nil, errors.New(fmt.Sprintf("syntax error: cannot find ("))
	}
	p.index++

	// 文字列
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.STRING {
		return nil, errors.New(fmt.Sprintf("syntax error: cannot find string"))
	}

	fileName := p.tokens[p.index].Content
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to open %s", fileName))
	}

	var jsonData map[string]interface{}
	if err := json.Unmarshal(bytes, &jsonData); err != nil {
		return nil, errors.New(fmt.Sprintf("failed to unmarshal %s", fileName))
	}

	p.index++

	// )
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.RPAREN {
		return nil, errors.New(fmt.Sprintf("syntax error: cannot find )"))
	}
	p.index++

	return jsonData, nil
}

// mapキー
func (p *Parser) mapKey() (values.MapKey, error) {
	// 変数名
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.IDENTIFIER {
		return values.MapKey{}, errors.New(fmt.Sprintf("syntax error: cannot find identifier"))
	}

	name := p.tokens[p.index].Content
	p.index++

	// .
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.DOT {
		return values.MapKey{}, errors.New(fmt.Sprintf("syntax error: cannot find ."))
	}

	p.index++

	// keys
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.KEYS {
		return values.MapKey{}, errors.New(fmt.Sprintf("syntax error: cannot find keys"))
	}

	p.index++

	return values.MapKey{Kind: values.MAPKEY, Name: name}, nil
}

// mapバリュー
func (p *Parser) mapValue() (values.MapValue, error) {
	// 変数名
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.IDENTIFIER {
		return values.MapValue{}, errors.New(fmt.Sprintf("syntax error: cannot find identifier"))
	}

	name := p.tokens[p.index].Content
	p.index++

	// [
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.LBRACKET {
		return values.MapValue{}, errors.New(fmt.Sprintf("syntax error: cannot find left bracket"))
	}

	p.index++

	v, err := p.singleAssignFormula()
	if err != nil {
		return values.MapValue{}, err
	}

	// [
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.RBRACKET {
		return values.MapValue{}, errors.New(fmt.Sprintf("syntax error: cannot find right bracket"))
	}

	p.index++

	return values.MapValue{Kind: values.MAPVALUE, Name: name, Key: v}, nil
}

// if節
// ifコードとdescriptionコードを分けて返す
// 間にJUMP命令が挟まるため、ifBlockで結合処理を行う
func (p *Parser) ifSection() ([]codes.Code, error) {
	var ifCodes []codes.Code
	var err error

	// "if"
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.IF {
		return nil, errors.New(fmt.Sprintf("syntax error: cannot find 'if'"))
	}

	p.index++

	// "("
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.LPAREN {
		return nil, errors.New(fmt.Sprintf("syntax error: cannot find '('"))
	}

	p.index++

	// 条件判定式
	condition, err := p.conditionalFormula()
	if err != nil {
		return nil, err
	}

	ifCode := codes.If{Kind: codes.IF, Condition: *condition}
	ifCodes = append(ifCodes, ifCode)

	// ")"
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.RPAREN {
		return nil, errors.New(fmt.Sprintf("syntax error: cannot find ')'"))
	}

	p.index++

	// "{"
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.LBRACE {
		return nil, errors.New(fmt.Sprintf("syntax error: cannot find '{'"))
	}

	p.index++

	// 記述部
	descCodes, err := p.description()
	if err != nil {
		return nil, err
	}

	ifCodes = append(ifCodes, descCodes...)

	// "}"
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.RBRACE {
		return nil, errors.New(fmt.Sprintf("syntax error: cannot find '}'"))
	}

	p.index++

	endCode := codes.End{Kind: codes.END}
	ifCodes = append(ifCodes, endCode)

	return ifCodes, nil
}

// elif節
// elifコードとdescriptionコードを分けて返す
// 間にJUMP命令が挟まるため、ifBlockで結合処理を行う
func (p *Parser) elifSection() ([]codes.Code, error) {
	var elifCodes []codes.Code
	var err error

	// "else if"
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.ELIF {
		return nil, errors.New(fmt.Sprintf("syntax error: cannot find 'else if'"))
	}

	p.index++

	// "("
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.LPAREN {
		return nil, errors.New(fmt.Sprintf("syntax error: cannot find '('"))
	}

	p.index++

	// 条件判定式
	condition, err := p.conditionalFormula()
	if err != nil {
		return nil, err
	}

	elifCode := codes.Elif{Kind: codes.ELIF, Condition: *condition}
	elifCodes = append(elifCodes, elifCode)

	// ")"
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.RPAREN {
		return nil, errors.New(fmt.Sprintf("syntax error: cannot find ')'"))
	}

	p.index++

	// "{"
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.LBRACE {
		return nil, errors.New(fmt.Sprintf("syntax error: cannot find '{'"))
	}

	p.index++

	// 記述部
	descCodes, err := p.description()
	if err != nil {
		return nil, err
	}

	elifCodes = append(elifCodes, descCodes...)

	// "}"
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.RBRACE {
		return nil, errors.New(fmt.Sprintf("syntax error: cannot find '}'"))
	}

	p.index++

	endCode := codes.End{Kind: codes.END}
	elifCodes = append(elifCodes, endCode)

	return elifCodes, nil
}

// else節
// elseの場合はJUMPなどないので全てコードをまとめて返す
func (p *Parser) elseSection() ([]codes.Code, error) {
	var elseCodes []codes.Code
	var err error

	// "else"
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.ELSE {
		return nil, errors.New(fmt.Sprintf("syntax error: cannot find 'else'"))
	}

	elseCode := codes.Else{Kind: codes.ELSE}
	elseCodes = append(elseCodes, elseCode)

	p.index++

	// "{"
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.LBRACE {
		return nil, errors.New(fmt.Sprintf("syntax error: cannot find '{'"))
	}

	p.index++

	// 記述部
	descCodes, err := p.description()
	if err != nil {
		return nil, err
	}

	elseCodes = append(elseCodes, descCodes...)

	// "}"
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.RBRACE {
		return nil, errors.New(fmt.Sprintf("syntax error: cannot find '}'"))
	}

	p.index++

	endCode := codes.End{Kind: codes.END}
	elseCodes = append(elseCodes, endCode)

	return elseCodes, nil
}

// forブロック
func (p *Parser) forBlock() ([]codes.Code, error) {
	var forCodes []codes.Code
	// "for"
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.FOR {
		return nil, errors.New(fmt.Sprintf("syntax error: cannot find 'for'"))
	}

	p.index++

	// "("
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.LPAREN {
		return nil, errors.New(fmt.Sprintf("syntax error: cannot find '('"))
	}

	p.index++

	// 変数名
	itrName, err := p.variableName()
	if err != nil {
		return nil, err
	}

	// "in"
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.IN {
		return nil, errors.New(fmt.Sprintf("syntax errir: cannot find 'in'"))
	}

	p.index++

	var value values.Value
	if p.index+1 < len(p.tokens) && p.tokens[p.index].Kind == token.IDENTIFIER && p.tokens[p.index+1].Kind == token.DOT {
		value, err = p.mapKey()
		if err != nil {
			return nil, err
		}
	} else if p.index+1 < len(p.tokens) && p.tokens[p.index].Kind == token.IDENTIFIER && p.tokens[p.index+1].Kind == token.LBRACKET {
		value, err = p.mapValue()
		if err != nil {
			return nil, err
		}
	} else {
		vName, err := p.variableName()
		if err != nil {
			return nil, err
		}
		value = values.Ident{Kind: values.IDENT, Name: vName}
	}

	// ")"
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.RPAREN {
		return nil, errors.New(fmt.Sprintf("syntax error: cannot find ')'"))
	}

	p.index++

	forCode := codes.For{Kind: codes.FOR, ItrName: itrName, ArrayValue: value}
	forCodes = append(forCodes, forCode)

	// "{"
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.LBRACE {
		return nil, errors.New(fmt.Sprintf("syntax error: cannot find '{'"))
	}

	p.index++

	// 記述部
	descCodes, err := p.description()
	if err != nil {
		return nil, err
	}

	forCodes = append(forCodes, descCodes...)

	// "}"
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.RBRACE {
		return nil, errors.New(fmt.Sprintf("syntax error: cannot find '}'"))
	}

	p.index++

	endCode := codes.End{Kind: codes.END}
	forCodes = append(forCodes, endCode)

	return forCodes, nil
}

// 関数名
func (p *Parser) functionName() (string, error) {
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.IDENTIFIER {
		return "", errors.New(fmt.Sprintf("syntax error: cannot find an identifier"))
	}

	funcName := p.tokens[p.index].Content

	p.index++

	return funcName, nil
}

// 変数名
func (p *Parser) variableName() (string, error) {
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.IDENTIFIER {
		return "", errors.New(fmt.Sprintf("syntax error: cannot find an identifier"))
	}

	varName := p.tokens[p.index].Content

	p.index++

	return varName, nil
}
