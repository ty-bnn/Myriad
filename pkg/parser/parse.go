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
	fmt.Printf("Parsing %s ...\n", p.filePath)

	err := p.program()
	if err != nil {
		return err
	}

	fmt.Printf("Parse %s Done.\n", p.filePath)
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

	t := tokenizer.NewTokenizer(lines, filePath)
	err = t.Tokenize()
	if err != nil {
		return err
	}

	newP := NewParser(t.Tokens, filePath)
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

	// 記述ブロック群
	funcCodes, err := p.descriptionBlockGroup()
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
	funcCodes, err := p.descriptionBlockGroup()
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

// main記述ブロック群
func (p *Parser) descriptionBlockGroup() ([]codes.Code, error) {
	var descCodes []codes.Code

	// "{"
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.LBRACE {
		return nil, errors.New(fmt.Sprintf("syntax error: cannot find '{'"))
	}

	p.index++

	for p.tokenIs(token.DFBEGIN, 0) || p.tokenIs(token.DFARG, 0) || p.tokenIs(token.IDENTIFIER, 0) || p.tokenIs(token.IF, 0) || p.tokenIs(token.FOR, 0) {
		if p.tokenIs(token.IDENTIFIER, 0) && p.tokenIs(token.DOUBLELESS, 1) {
			outCodes, err := p.outputBlock()
			if err != nil {
				return nil, err
			}
			descCodes = append(descCodes, outCodes...)
			continue
		}

		descBCodes, err := p.descriptionBlock()
		if err != nil {
			return nil, err
		}

		descCodes = append(descCodes, descBCodes...)
	}

	// "}"
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.RBRACE {
		return nil, errors.New(fmt.Sprintf("syntax error: cannot find '}'"))
	}

	p.index++

	return descCodes, nil
}

// 出力ブロック
func (p *Parser) outputBlock() ([]codes.Code, error) {
	var codeBlock []codes.Code

	vName, err := p.variableName()
	if err != nil {
		return nil, err
	}
	codeBlock = append(codeBlock, codes.Output{
		Kind: codes.OUTPUT,
		FilePath: values.Ident{
			Kind: values.IDENT,
			Name: vName,
		},
	})

	// "<<"
	if !p.tokenIs(token.DOUBLELESS, 0) {
		return nil, errors.New(fmt.Sprintf("syntax error: cannot find '<<'"))
	}
	p.index++

	descCodes, err := p.descriptionBlockGroup()
	if err != nil {
		return nil, err
	}
	codeBlock = append(codeBlock, descCodes...)

	codeBlock = append(codeBlock, codes.End{
		Kind: codes.END,
	})

	return codeBlock, nil
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
	} else if p.tokenIs(token.IDENTIFIER, 0) && p.tokenIs(token.DOT, 1) && p.tokenIs(token.APPEND, 2) {
		appendCode, err := p.appendArray()
		if err != nil {
			return nil, err
		}
		return []codes.Code{appendCode}, nil
	} else if p.tokenIs(token.IDENTIFIER, 0) && p.tokenIs(token.DOT, 1) && p.tokenIs(token.SORT, 2) {
		sortCode, err := p.sortArray()
		if err != nil {
			return nil, err
		}
		return []codes.Code{sortCode}, nil
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
	var assignValues []values.Value
	var err error
	// 代入値
	value, err := p.assignValue()
	if err != nil {
		return nil, err
	}
	assignValues = append(assignValues, value)

	for {
		// ","
		if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.COMMA {
			break
		}
		p.index++
		// 代入値
		value, err := p.assignValue()
		if err != nil {
			return assignValues, err
		}
		assignValues = append(assignValues, value)
	}

	return assignValues, nil
}

// ifブロック
func (p *Parser) ifBlock() ([]codes.Code, error) {
	var ifBCodes []codes.Code
	var err error
	type Jump struct {
		ptr int
		codes.Jump
	}
	var jumps []Jump

	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.IF {
		return nil, errors.New(fmt.Sprintf("syntax error: cannot find 'if'"))
	}

	ifCodes, err := p.ifSection()
	if err != nil {
		return nil, err
	}
	ifBCodes = append(ifBCodes, ifCodes...)
	jumps = append(jumps, Jump{ptr: len(ifBCodes) - len(ifCodes), Jump: codes.Jump{False: len(ifCodes)}})

	for {
		if !p.tokenIs(token.ELSE, 0) || !p.tokenIs(token.IF, 1) {
			break
		}

		// elif節
		elifCodes, err := p.elifSection()
		if err != nil {
			return nil, err
		}
		ifBCodes = append(ifBCodes, elifCodes...)
		jumps = append(jumps, Jump{ptr: len(ifBCodes) - len(elifCodes), Jump: codes.Jump{False: len(elifCodes)}})
	}

	if p.index < len(p.tokens) && p.tokens[p.index].Kind == token.ELSE {
		elseCodes, err := p.elseSection()
		if err != nil {
			return nil, err
		}
		ifBCodes = append(ifBCodes, elseCodes...)
	}

	length := len(ifBCodes)
	for i := 0; i < len(jumps); i++ {
		jumps[i].True = length
		length -= jumps[i].False
	}

	for _, jump := range jumps {
		switch ifBCodes[jump.ptr].GetKind() {
		case codes.IF:
			code := ifBCodes[jump.ptr].(codes.If)
			code.Jump = jump.Jump
			ifBCodes[jump.ptr] = code
		case codes.ELIF:
			code := ifBCodes[jump.ptr].(codes.Elif)
			code.Jump = jump.Jump
			ifBCodes[jump.ptr] = code
		}
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

	stackIndex := p.index
	compFml, err := p.compFormula()
	if err == nil {
		return compFml, nil
	}

	p.index = stackIndex
	compFml, err = p.analyzeStringFormula()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("syntax error: cannot find compFormula or analyzeStringFormula"))
	}
	return compFml, nil
}

// 比較式
func (p *Parser) compFormula() (*codes.ConditionalNode, error) {
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

// 文字列解析式
func (p *Parser) analyzeStringFormula() (*codes.ConditionalNode, error) {
	var falseFlag bool
	if p.tokenIs(token.NOT, 0) {
		falseFlag = true
		p.index++
	}

	left, err := p.singleAssignValue()
	if err != nil {
		return nil, err
	}
	lNode := codes.ConditionalNode{Var: left}

	if !p.tokenIs(token.DOT, 0) {
		return nil, errors.New(fmt.Sprintf("syntax error: cannot find '.'"))
	}
	p.index++

	var op codes.OperatorKind
	if p.tokenIs(token.STARTWITH, 0) {
		op = codes.STARTWITH
	} else if p.tokenIs(token.ENDWITH, 0) {
		op = codes.ENDWITH
	} else {
		return nil, errors.New(fmt.Sprintf("syntax error: cannot find 'startwith' or 'endwith'"))
	}
	p.index++

	if !p.tokenIs(token.LPAREN, 0) {
		return nil, errors.New(fmt.Sprintf("syntax error: cannot find '('"))
	}
	p.index++

	right, err := p.singleAssignValue()
	if err != nil {
		return nil, err
	}
	rNode := codes.ConditionalNode{Var: right}

	if !p.tokenIs(token.RPAREN, 0) {
		return nil, errors.New(fmt.Sprintf("syntax error: cannot find ')'"))
	}
	p.index++

	return &codes.ConditionalNode{Operator: op, Left: &lNode, Right: &rNode, False: falseFlag}, nil
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
	stackIndex := p.index
	complexValue, err := p.complexAssignValue()
	if err == nil {
		return complexValue, nil
	}
	p.index = stackIndex

	singleValue, err := p.singleAssignFormula()
	if err == nil {
		return singleValue, nil
	}

	return nil, errors.New(fmt.Sprintf("syntax error: cannot find assign value"))
}

// 単一代入式
func (p *Parser) singleAssignFormula() (values.Value, error) {
	var (
		value values.Value
		err   error
	)

	stackIndex := p.index
	value, err = p.trimStringFormula()
	if err != nil {
		p.index = stackIndex
		value, err = p.singleAssignValue()
		if err != nil {
			return value, err
		}
	}

	if !p.tokenIs(token.PLUS, 0) {
		return value, nil
	}

	var vls []values.Value
	vls = append(vls, value)

	for p.tokenIs(token.PLUS, 0) {
		// +
		p.index++

		stackIndex = p.index
		value, err = p.trimStringFormula()
		if err != nil {
			p.index = stackIndex
			value, err = p.singleAssignValue()
			if err != nil {
				return value, err
			}
		}
		vls = append(vls, value)
	}

	return values.AddString{Kind: values.ADDSTRING, Values: vls}, nil
}

// 文字列除去式
func (p *Parser) trimStringFormula() (values.TrimString, error) {
	target, err := p.singleAssignValue()
	if err != nil {
		return values.TrimString{}, err
	}

	if !p.tokenIs(token.DOT, 0) {
		return values.TrimString{}, errors.New(fmt.Sprintf("syntax error: cannot find '.'"))
	}
	p.index++

	var from values.FromKind
	if p.tokenIs(token.TRIMLEFT, 0) {
		from = values.LEFT
	} else if p.tokenIs(token.TRIMRIGHT, 0) {
		from = values.RIGHT
	} else {
		return values.TrimString{}, errors.New(fmt.Sprintf("syntax error: cannot find 'leftTrim' or 'rightTrim'"))
	}
	p.index++

	if !p.tokenIs(token.LPAREN, 0) {
		return values.TrimString{}, errors.New(fmt.Sprintf("syntax error: cannot find '('"))
	}
	p.index++

	trim, err := p.singleAssignValue()
	if err != nil {
		return values.TrimString{}, err
	}

	if !p.tokenIs(token.RPAREN, 0) {
		return values.TrimString{}, errors.New(fmt.Sprintf("syntax error: cannot find ')'"))
	}
	p.index++

	return values.TrimString{Kind: values.TRIMSTRING, Target: target, Trim: trim, From: from}, nil
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
	stackIndex := p.index
	arrValues, err := p.array()
	if err == nil {
		return values.Literals{Kind: values.LITERALS, Values: arrValues}, nil
	}
	p.index = stackIndex
	jsonData, err := p.jsonUnmarshal()
	if err == nil {
		return values.Map{Kind: values.MAP, Value: jsonData}, nil
	}
	p.index = stackIndex
	mapValue, err := p.mapKey()
	if err == nil {
		return mapValue, nil
	}
	p.index = stackIndex
	splitArr, err := p.splitStringFormula()
	if err == nil {
		return splitArr, nil
	}
	return nil, errors.New(fmt.Sprintf("syntax error: cannot parse complex assign value"))
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

	var keys []values.Value
	for {
		if !p.tokenIs(token.LBRACKET, 0) {
			break
		}

		// [
		if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.LBRACKET {
			return values.MapValue{}, errors.New(fmt.Sprintf("syntax error: cannot find left bracket"))
		}

		p.index++

		key, err := p.singleAssignFormula()
		if err != nil {
			return values.MapValue{}, err
		}
		keys = append(keys, key)

		// ]
		if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.RBRACKET {
			return values.MapValue{}, errors.New(fmt.Sprintf("syntax error: cannot find right bracket"))
		}

		p.index++
	}

	return values.MapValue{Kind: values.MAPVALUE, Name: name, Keys: keys}, nil
}

// 文字列分割式
func (p *Parser) splitStringFormula() (values.SplitString, error) {
	target, err := p.singleAssignValue()
	if err != nil {
		return values.SplitString{}, err
	}

	if !p.tokenIs(token.DOT, 0) {
		return values.SplitString{}, errors.New(fmt.Sprintf("syntax error: cannot find '.'"))
	}
	p.index++

	if !p.tokenIs(token.SPLIT, 0) {
		return values.SplitString{}, errors.New(fmt.Sprintf("syntax error: cannot find 'split'"))
	}
	p.index++

	if !p.tokenIs(token.LPAREN, 0) {
		return values.SplitString{}, errors.New(fmt.Sprintf("syntax error: cannot find '('"))
	}
	p.index++

	sep, err := p.singleAssignValue()
	if err != nil {
		return values.SplitString{}, err
	}

	if !p.tokenIs(token.RPAREN, 0) {
		return values.SplitString{}, errors.New(fmt.Sprintf("syntax error: cannot find ')'"))
	}
	p.index++

	return values.SplitString{Kind: values.SPLITSTRING, Target: target, Sep: sep}, nil
}

// if節
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

	// 記述ブロック群
	descCodes, err := p.descriptionBlockGroup()
	if err != nil {
		return nil, err
	}

	ifCodes = append(ifCodes, descCodes...)

	endCode := codes.End{Kind: codes.END}
	ifCodes = append(ifCodes, endCode)

	return ifCodes, nil
}

// elif節
func (p *Parser) elifSection() ([]codes.Code, error) {
	var elifCodes []codes.Code
	var err error

	// "else"
	if !p.tokenIs(token.ELSE, 0) {
		return nil, errors.New(fmt.Sprintf("syntax error: cannot find 'else'"))
	}
	p.index++

	// "if"
	if !p.tokenIs(token.IF, 0) {
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

	elifCode := codes.Elif{Kind: codes.ELIF, Condition: *condition}
	elifCodes = append(elifCodes, elifCode)

	// ")"
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.RPAREN {
		return nil, errors.New(fmt.Sprintf("syntax error: cannot find ')'"))
	}
	p.index++

	// 記述ブロック群
	descCodes, err := p.descriptionBlockGroup()
	if err != nil {
		return nil, err
	}

	elifCodes = append(elifCodes, descCodes...)

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

	// 記述ブロック群
	descCodes, err := p.descriptionBlockGroup()
	if err != nil {
		return nil, err
	}

	elseCodes = append(elseCodes, descCodes...)

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

	// 記述ブロック群
	descCodes, err := p.descriptionBlockGroup()
	if err != nil {
		return nil, err
	}

	forCodes = append(forCodes, descCodes...)

	endCode := codes.End{Kind: codes.END}
	forCodes = append(forCodes, endCode)

	return forCodes, nil
}

// 配列登録式
func (p *Parser) appendArray() (codes.Append, error) {
	arrayName, err := p.variableName()
	if err != nil {
		return codes.Append{}, err
	}
	if !p.tokenIs(token.DOT, 0) {
		return codes.Append{}, errors.New(fmt.Sprintf("syntax error: cannot find '.'"))
	}
	p.index++
	if !p.tokenIs(token.APPEND, 0) {
		return codes.Append{}, errors.New(fmt.Sprintf("syntax error: cannot find 'append'"))
	}
	p.index++
	if !p.tokenIs(token.LPAREN, 0) {
		return codes.Append{}, errors.New(fmt.Sprintf("syntax error: cannot find '('"))
	}
	p.index++
	elem, err := p.singleAssignFormula()
	if err != nil {
		return codes.Append{}, err
	}
	if !p.tokenIs(token.RPAREN, 0) {
		return codes.Append{}, errors.New(fmt.Sprintf("syntax error: cannot find ')'"))
	}
	p.index++
	return codes.Append{Kind: codes.APPEND, Array: arrayName, Element: elem}, nil
}

// 配列ソート文
func (p *Parser) sortArray() (codes.Sort, error) {
	arrayName, err := p.variableName()
	if err != nil {
		return codes.Sort{}, err
	}
	if !p.tokenIs(token.DOT, 0) {
		return codes.Sort{}, errors.New(fmt.Sprintf("syntax error: cannot find '.'"))
	}
	p.index++
	if !p.tokenIs(token.SORT, 0) {
		return codes.Sort{}, errors.New(fmt.Sprintf("syntax error: cannot find 'sort'"))
	}
	p.index++
	if !p.tokenIs(token.LPAREN, 0) {
		return codes.Sort{}, errors.New(fmt.Sprintf("syntax error: cannot find '('"))
	}
	p.index++
	if !p.tokenIs(token.RPAREN, 0) {
		return codes.Sort{}, errors.New(fmt.Sprintf("syntax error: cannot find ')'"))
	}
	p.index++
	return codes.Sort{Kind: codes.SORT, Array: arrayName}, nil
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
