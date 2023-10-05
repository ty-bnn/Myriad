package parser

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ty-bnn/myriad/pkg/model/codes"
	"github.com/ty-bnn/myriad/pkg/model/token"
	"github.com/ty-bnn/myriad/pkg/model/vars"
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

	// 変数
	argName, err := p.variable()
	if err != nil {
		return nil, err
	}

	defCode := codes.Define{
		Kind: codes.DEFINE,
		Var:  vars.Single{Kind: vars.SINGLE, Name: argName},
	}

	argCodes = append(argCodes, defCode)

	for {
		// ","
		if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.COMMA {
			break
		}

		p.index++

		// 変数
		argName, err = p.variable()
		if err != nil {
			return nil, err
		}

		defCode = codes.Define{
			Kind: codes.DEFINE,
			Var:  vars.Single{Kind: vars.SINGLE, Name: argName},
		}

		argCodes = append(argCodes, defCode)
	}

	return argCodes, nil
}

// 変数
// TODO: 変数いらない、変数名だけで良い
func (p *Parser) variable() (string, error) {
	var err error

	// 変数名
	varName, err := p.variableName()
	if err != nil {
		return "", err
	}

	return varName, nil
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
		if p.index >= len(p.tokens) || (p.tokens[p.index].Kind != token.DFCOMMAND && p.tokens[p.index].Kind != token.DFARG && p.tokens[p.index].Kind != token.IDENTIFIER && p.tokens[p.index].Kind != token.IF && p.tokens[p.index].Kind != token.FOR) {
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
	if p.index < len(p.tokens) && p.tokens[p.index].Kind == token.DFCOMMAND || p.tokens[p.index].Kind == token.DFARG {
		// Dfile文
		dockerCodes, err := p.dockerFile()
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

// Dfile文
func (p *Parser) dockerFile() ([]codes.Code, error) {
	var dockerCodes []codes.Code
	var err error
	// Df命令
	if p.index >= len(p.tokens) || (p.tokens[p.index].Kind != token.DFCOMMAND && p.tokens[p.index].Kind != token.DFARG && p.tokens[p.index].Kind != token.LDOUBLEBRA) {
		return nil, errors.New(fmt.Sprintf("syntax error: cannot find a Dockerfile comamnd"))
	}

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

	// 置換変数
	target, err := p.formula()
	if err != nil {
		return nil, err
	}

	repCode := codes.Replace{Kind: codes.REPLACE, RepVar: target}

	// }}
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.RDOUBLEBRA {
		return nil, errors.New(fmt.Sprintf("syntax error: cannot find '}}'"))
	}

	p.index++

	return repCode, nil
}

// 関数呼び出し文
func (p *Parser) functionCall() (codes.Code, error) {
	var args []vars.Var
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

	// 式の並び
	if p.index < len(p.tokens) && (p.tokens[p.index].Kind == token.STRING || p.tokens[p.index].Kind == token.IDENTIFIER) {
		args, err = p.rowOfFormulas()
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

// 式の並び
func (p *Parser) rowOfFormulas() ([]vars.Var, error) {
	var fmls []vars.Var
	var err error
	// 式
	fml, err := p.formula()
	if err != nil {
		return fmls, err
	}

	fmls = append(fmls, fml)

	for {
		// ","
		if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.COMMA {
			break
		}

		p.index++

		// 式
		fml, err := p.formula()
		if err != nil {
			return fmls, err
		}

		fmls = append(fmls, fml)
	}

	return fmls, nil
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
func (p *Parser) conditionalFormula() (codes.Condition, error) {
	var err error
	// 式
	left, err := p.formula()
	if err != nil {
		return codes.Condition{}, err
	}

	// 比較演算子
	op, err := p.conditionalOperator()
	if err != nil {
		return codes.Condition{}, err
	}

	// 式
	right, err := p.formula()
	if err != nil {
		return codes.Condition{}, err
	}

	return codes.Condition{Left: left, Right: right, Operator: op}, nil
}

// 式
func (p *Parser) formula() (vars.Var, error) {
	var fml vars.Var

	// 文字列でも変数でもない場合
	if p.index >= len(p.tokens) || (p.tokens[p.index].Kind != token.STRING && p.tokens[p.index].Kind != token.IDENTIFIER) {
		return fml, errors.New(fmt.Sprintf("syntax error: cannot find a formula"))
	}

	// 文字列
	if p.index < len(p.tokens) && p.tokens[p.index].Kind == token.STRING {
		fml = vars.Literal{Kind: vars.LITERAL, Value: p.tokens[p.index].Content}
		p.index++
		return fml, nil
	}

	// 変数 | 配列要素
	name := p.tokens[p.index].Content
	p.index++

	// [
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.LBRACKET {
		return vars.Single{Kind: vars.SINGLE, Name: name}, nil
	}
	p.index++

	// 数字
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.NUMBER {
		return fml, errors.New(fmt.Sprintf("syntax error: cannot find a number"))
	}

	index, err := strconv.Atoi(p.tokens[p.index].Content)
	if err != nil {
		return fml, err
	}

	p.index++

	// ]
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.RBRACKET {
		return fml, errors.New(fmt.Sprintf("syntax error: cannot find a right bracket"))
	}
	p.index++

	return vars.Element{Kind: vars.ELEMENT, Name: name, Index: index}, nil
}

// 比較演算子
func (p *Parser) conditionalOperator() (codes.OpeKind, error) {
	var op codes.OpeKind

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

	variable, err := p.assignValue(vName)
	if err != nil {
		return nil, err
	}

	defCode := codes.Define{Kind: codes.DEFINE, Var: variable}

	return defCode, nil
}

// 変数代入文
func (p *Parser) assignVariable() (codes.Code, error) {
	var err error

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

	variable, err := p.assignValue(vName)
	if err != nil {
		return nil, err
	}

	defCode := codes.Assign{Kind: codes.ASSIGN, Var: variable}

	return defCode, nil
}

// 代入値
func (p *Parser) assignValue(vName string) (vars.Var, error) {
	var variable vars.Var
	// 文字列
	if p.index < len(p.tokens) && p.tokens[p.index].Kind == token.STRING {
		variable = vars.Single{Kind: vars.SINGLE, Name: vName, Value: p.tokens[p.index].Content}

		p.index++
		// 配列
	} else if p.index < len(p.tokens) && p.tokens[p.index].Kind == token.LBRACE {
		array, err := p.array()
		if err != nil {
			return nil, err
		}

		variable = vars.Array{Kind: vars.ARRAY, Name: vName, Values: array}
	} else {
		return nil, errors.New(fmt.Sprintf("syntax error: cannot find string or '{'"))
	}

	return variable, nil
}

// 配列
func (p *Parser) array() ([]string, error) {
	var err error

	// {
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.LBRACE {
		return []string{}, errors.New(fmt.Sprintf("syntax error: cannot find '{'"))
	}

	p.index++

	array, err := p.rowOfStrings()
	if err != nil {
		return array, err
	}

	// }
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.RBRACE {
		return nil, errors.New(fmt.Sprintf("syntax error: cannot find '}'"))
	}

	p.index++

	return array, nil
}

// 文字列の並び
func (p *Parser) rowOfStrings() ([]string, error) {
	var strings []string

	// 文字列
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.STRING {
		return strings, errors.New(fmt.Sprintf("syntax error: cannot find string"))
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

	ifCode := codes.If{Kind: codes.IF, Condition: condition}
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

	elifCode := codes.Elif{Kind: codes.ELIF, Condition: condition}
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

	itrVar := vars.Single{Kind: vars.SINGLE, Name: itrName}

	// "in"
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.IN {
		return nil, errors.New(fmt.Sprintf("syntax errir: cannot find 'in'"))
	}

	p.index++

	// 変数名
	arrayName, err := p.variableName()
	if err != nil {
		return nil, err
	}

	arrayVar := vars.Array{Kind: vars.ARRAY, Name: arrayName}

	// ")"
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != token.RPAREN {
		return nil, errors.New(fmt.Sprintf("syntax error: cannot find ')'"))
	}

	p.index++

	forCode := codes.For{Kind: codes.FOR, Itr: itrVar, Array: arrayVar}
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
