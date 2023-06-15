package parser

import (
	"fmt"
	"errors"

	"myriad/tokenizer"
)

func (p *Parser) program() error {
	var err error

	// { 関数インポート文 }
	for ;; {
		if p.index >= len(p.tokens) || p.tokens[p.index].Kind != tokenizer.SIMPORT {
			break
		}

		err = p.importFunc()
		if err != nil {
			return err
		}
	}

	// { 関数 }
	for ;; {
		if p.index >= len(p.tokens) || p.tokens[p.index].Kind != tokenizer.SIDENTIFIER {
			break
		}

		err = p.function()
		if err != nil {
			return err
		}
	}

	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != tokenizer.SMAIN {
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
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != tokenizer.SIMPORT {
		return errors.New(fmt.Sprintf("syntax error: cannot find 'import'"))
	}

	p.index++

	// 関数名
	err = p.functionName()
	if err != nil {
		return err
	}

	// "from"
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != tokenizer.SFROM {
		return errors.New(fmt.Sprintf("syntax error: cannot find 'from'"))
	}

	p.index++

	// ファイル名
	err = p.fileName()
	if err != nil {
		return err
	}

	return nil
}

// ファイル名
func (p *Parser) fileName() error {
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != tokenizer.SSTRING {
		return errors.New(fmt.Sprintf("syntax error: cannot find an identifier"))
	}

	p.index++

	return nil
}

// 関数
func (p *Parser) function() error {
	var err error

	// 関数名
	err = p.functionName()
	if err != nil {
		return err
	}

	// 引数宣言部
	err = p.argumentDecralation()
	if err != nil {
		return err
	}

	// 関数記述部
	err = p.functionDescription()
	if err != nil {
		return err
	}

	return nil
}

// メイン部
func (p *Parser) mainFunction() error {
	var err error

	// "main"
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != tokenizer.SMAIN {
		return errors.New(fmt.Sprintf("syntax error: cannot find 'main' to declare"))
	}

	p.index++

	// 引数宣言部
	err = p.argumentDecralation()
	if err != nil {
		return err
	}

	// 関数記述部
	err = p.functionDescription()
	if err != nil {
		return err
	}

	return nil
}

// 引数宣言部
func (p *Parser) argumentDecralation() error {
	var err error

	// "("
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != tokenizer.SLPAREN {
		return errors.New(fmt.Sprintf("syntax error: cannot find '('"))
	}

	p.index++

	// 引数群
	if 	p.index < len(p.tokens) && p.tokens[p.index].Kind == tokenizer.SIDENTIFIER {
		err = p.arguments()
		if err != nil {
			return err
		}
	}

	// ")"
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != tokenizer.SRPAREN {
		return errors.New(fmt.Sprintf("syntax error: cannot find ')'"))
	}

	p.index++

	return nil
}

// 引数群
func (p *Parser) arguments() error {
	var err error

	// 変数
	err = p.variable()
	if err != nil {
		return err
	}

	for ;; {
		// ","
		if p.index >= len(p.tokens) || p.tokens[p.index].Kind != tokenizer.SCOMMA {
			break
		}

		p.index++
		
		// 変数
		err = p.variable()
		if err != nil {
			return err
		}
	}

	return nil
}

// 変数
func (p *Parser) variable() error {
	var err error

	// 変数名
	err = p.variableName()
	if err != nil {
		return err
	}

	// "[]"
	if p.index < len(p.tokens) && p.tokens[p.index].Kind == tokenizer.SARRANGE {
		p.index++
	}

	return nil
}

// 関数記述部
func (p *Parser) functionDescription() error {
	var err error

	// "{"
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != tokenizer.SLBRACE {
		return errors.New(fmt.Sprintf("syntax error: cannot find '{'"))
	}

	p.index++

	// 記述部
	err = p.description()
	if err != nil {
		return err
	}

	// "}"
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != tokenizer.SRBRACE {
		return errors.New(fmt.Sprintf("syntax error: cannot find '}'"))
	}

	p.index++

	return nil
}

// 記述部
func (p *Parser) description() error {
	var err error

	// 記述ブロック
	err = p.descriptionBlock()
	if err != nil {
		return err
	}

	for ;; {
		if p.index >= len(p.tokens) || (p.tokens[p.index].Kind != tokenizer.SDFCOMMAND && p.tokens[p.index].Kind != tokenizer.SDFARG && p.tokens[p.index].Kind != tokenizer.SIDENTIFIER && p.tokens[p.index].Kind != tokenizer.SIF) {
			break
		}
			
		err = p.descriptionBlock()
		if err != nil {
			return err
		}
	}

	return nil
}

// 記述ブロック
func (p *Parser) descriptionBlock() error {
	var err error

	if p.index < len(p.tokens) && p.tokens[p.index].Kind == tokenizer.SDFCOMMAND || p.tokens[p.index].Kind == tokenizer.SDFARG {
		// Dfile文
		err = p.dockerFile()
		if err != nil {
			return err
		}
	} else if p.index + 1 < len(p.tokens) && p.tokens[p.index].Kind == tokenizer.SIDENTIFIER && p.tokens[p.index + 1].Kind == tokenizer.SLPAREN {
		// 関数呼び出し文
		err = p.functionCall()
		if err != nil {
			return err
		}
	} else if p.index + 1 < len(p.tokens) && p.tokens[p.index].Kind == tokenizer.SIDENTIFIER && p.tokens[p.index + 1].Kind == tokenizer.SDEFINE {
		// 変数定義文
		err = p.defineVariable()
	} else if p.index < len(p.tokens) && p.tokens[p.index].Kind == tokenizer.SIF {
		// ifブロック
		err = p.ifBlock()
		if err != nil {
			return err
		}
	} else {
		return errors.New(fmt.Sprintf("syntax error: cannot find a description block"))
	}

	return nil
}

// Dfile文
func (p *Parser) dockerFile() error {
	var err error
	// Df命令
	if p.index >= len(p.tokens) || (p.tokens[p.index].Kind != tokenizer.SDFCOMMAND && p.tokens[p.index].Kind != tokenizer.SDFARG && p.tokens[p.index].Kind != tokenizer.SLDOUBLEBRA ) {
		return errors.New(fmt.Sprintf("syntax error: cannot find a Dockerfile comamnd"))
	}

	if p.tokens[p.index].Kind == tokenizer.SDFCOMMAND {
		p.index++
	}


	// Df引数部
	err = p.dfArgs()
	if err != nil {
		return err
	}

	return nil
}

// Df引数部
func (p *Parser) dfArgs() error {
	var err error
	err = p.dfArg()
	if err != nil {
		return err
	}

	for ;; {
		if p.index < len(p.tokens) && (p.tokens[p.index].Kind == tokenizer.SDFARG || p.tokens[p.index].Kind == tokenizer.SLDOUBLEBRA) {
			err = p.dfArg()
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
func (p *Parser) dfArg() error {
	if p.index < len(p.tokens) && p.tokens[p.index].Kind == tokenizer.SDFARG {
		// 生のDf引数
		p.index++
	} else if p.index < len(p.tokens) && p.tokens[p.index].Kind == tokenizer.SLDOUBLEBRA {
		//置換式
		err := p.replaceFormula()
		if err != nil {
			return err
		}
	} else {
		return errors.New(fmt.Sprintf("syntax error: cannot find Df argument"))
	}

	return nil
}

// 置換式
func (p *Parser) replaceFormula() error {
	// {{
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != tokenizer.SLDOUBLEBRA {
		return errors.New(fmt.Sprintf("syntax error: cannot find '{'"))
	}
	
	p.index++

	// 置換変数
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != tokenizer.SASSIGNVARIABLE {
		return errors.New(fmt.Sprintf("syntax error: cannot find assign variable"))
	}

	p.index++

	// 配列要素
	if p.index < len(p.tokens) && p.tokens[p.index].Kind == tokenizer.SLBRACKET {
		// [
		p.index++

		// 数字
		if p.index >= len(p.tokens) || p.tokens[p.index].Kind != tokenizer.SNUMBER {
			return errors.New(fmt.Sprintf("syntax error: cannot find array index"))
		}

		p.index++

		// ]
		if p.index >= len(p.tokens) || p.tokens[p.index].Kind != tokenizer.SRBRACKET {
			return errors.New(fmt.Sprintf("syntax error: cannot find ']'"))
		}

		p.index++
	}

	// }}
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != tokenizer.SRDOUBLEBRA {
		return errors.New(fmt.Sprintf("syntax error: cannot find '}'"))
	}

	p.index++

	return nil
}

// 関数呼び出し文
func (p *Parser) functionCall() error {
	var err error

	// 関数名
	err = p.functionName()
	if err != nil {
		return err
	}

	// "("
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != tokenizer.SLPAREN {
		return errors.New(fmt.Sprintf("syntax error: cannot find '('"))
	}

	p.index++

	// 式の並び
	if p.index < len(p.tokens) && (p.tokens[p.index].Kind == tokenizer.SSTRING || p.tokens[p.index].Kind == tokenizer.SIDENTIFIER) {
		err = p.rowOfFormulas()
		if err != nil {
			return err
		}
	}

	// ")"
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != tokenizer.SRPAREN {
		return errors.New(fmt.Sprintf("syntax error: cannot find ')'"))
	}

	p.index++

	return nil
}

// 式の並び
func (p *Parser) rowOfFormulas() error {
	var err error
	// 式
	err = p.formula()
	if err != nil {
		return err
	}

	for ;; {
		// ","
		if p.index >= len(p.tokens) || p.tokens[p.index].Kind != tokenizer.SCOMMA {
			break
		}

		p.index++


		// 式
		err = p.formula()
		if err != nil {
			return err
		}
	}

	return nil
}

// ifブロック
func (p *Parser) ifBlock() error {
	var err error

	// "if"
	if p.index >= len(p.tokens) && p.tokens[p.index].Kind != tokenizer.SIF {
		return errors.New(fmt.Sprintf("syntax error: cannot find 'if'"))
	}

	p.index++

	// "("
	if p.index >= len(p.tokens) && p.tokens[p.index].Kind != tokenizer.SLPAREN {
		return errors.New(fmt.Sprintf("syntax error: cannot find '('"))
	}

	p.index++

	// 条件判定式
	err = p.conditionalFormula()
	if err != nil {
		return err
	}

	// ")"
	if p.index >= len(p.tokens) && p.tokens[p.index].Kind != tokenizer.SRPAREN {
		return errors.New(fmt.Sprintf("syntax error: cannot find ')'"))
	}

	p.index++

	// "{"
	if p.index >= len(p.tokens) && p.tokens[p.index].Kind != tokenizer.SLBRACE {
		return errors.New(fmt.Sprintf("syntax error: cannot find '{'"))
	}
	
	p.index++

	// 記述部
	err = p.description()
	if err != nil {
		return err
	}

	// "}"
	if p.index >= len(p.tokens) && p.tokens[p.index].Kind != tokenizer.SRBRACE {
		return errors.New(fmt.Sprintf("syntax error: cannot find '}'"))
	}

	p.index++

	for ;; {
		if p.index >= len(p.tokens) || p.tokens[p.index].Kind != tokenizer.SELIF {
			break
		}

		// elif節
		err = p.elifSection()
		if err != nil {
			return err
		}
	}

	if p.index < len(p.tokens) && p.tokens[p.index].Kind == tokenizer.SELSE {
		err = p.elseSection()
		if err != nil {
			return err
		}
	}

	return nil
}

// 条件判定式
func (p *Parser) conditionalFormula() error {
	var err error
	// 式
	err = p.formula()
	if err != nil {
		return err
	}

	// 比較演算子
	err = p.conditionalOperator()
	if err != nil {
		return err
	}

	// 式
	err = p.formula()
	if err != nil {
		return err
	}

	return nil
}

// 式
func (p *Parser) formula() error {
	// 変数, 文字列
	if p.index >= len(p.tokens) || (p.tokens[p.index].Kind != tokenizer.SIDENTIFIER && p.tokens[p.index].Kind != tokenizer.SSTRING) {
		return errors.New(fmt.Sprintf("syntax error: cannot find a formula"))
	}

	p.index++

	return nil
}

// 比較演算子
func (p *Parser) conditionalOperator() error {
	// "==", "!="
	if p.index >= len(p.tokens) || (p.tokens[p.index].Kind != tokenizer.SEQUAL && p.tokens[p.index].Kind != tokenizer.SNOTEQUAL) {
		return errors.New(fmt.Sprintf("syntax error: cannot find conditional operator"))
	}

	p.index++

	return nil
}

// 変数定義文
func (p *Parser) defineVariable() error {
	var err error

	// 変数名
	err = p.variableName()
	if err != nil {
		return err
	}

	// ":="
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != tokenizer.SDEFINE {
		return errors.New(fmt.Sprintf("syntax error: cannot find ':='"))
	}

	p.index++

	// 文字列
	if p.index < len(p.tokens) && p.tokens[p.index].Kind == tokenizer.SSTRING {
		p.index++
	// 配列
	} else if p.index < len(p.tokens) && p.tokens[p.index].Kind == tokenizer.SLBRACE {
		err = p.array()
		if err != nil {
			return err
		}
	} else {
		return errors.New(fmt.Sprintf("syntax error: cannot find string or '{'"))
	}

	return nil
}

// 配列
func (p *Parser) array() error {
	var err error

	// {
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != tokenizer.SLBRACE {
		return errors.New(fmt.Sprintf("syntax error: cannot find '{'"))
	}

	p.index++

	err = p.rowOfStrings()
	if err != nil {
		return err
	}

	// }
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != tokenizer.SRBRACE {
		return errors.New(fmt.Sprintf("syntax error: cannot find '}'"))
	}

	p.index++

	return nil
}

// 文字列の並び
func (p *Parser) rowOfStrings() error {
	// 文字列
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != tokenizer.SSTRING {
		return errors.New(fmt.Sprintf("syntax error: cannot find string"))
	}

	p.index++

	for ;; {
		// ","
		if p.index >= len(p.tokens) || p.tokens[p.index].Kind != tokenizer.SCOMMA {
			break
		}

		p.index++


		// 文字列
		if p.index >= len(p.tokens) || p.tokens[p.index].Kind != tokenizer.SSTRING {
			return errors.New(fmt.Sprintf("syntax error: cannot find string"))
		}

		p.index++
	}

	return nil
}

// elif節
func (p *Parser) elifSection() error {
	var err error

	// "else if"
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != tokenizer.SELIF {
		return errors.New(fmt.Sprintf("syntax error: cannot find 'else if'"))
	}

	p.index++

	// "("
	if p.index >= len(p.tokens) && p.tokens[p.index].Kind != tokenizer.SLPAREN {
		return errors.New(fmt.Sprintf("syntax error: cannot find '('"))
	}

	p.index++

	// 条件判定式
	err = p.conditionalFormula()
	if err != nil {
		return err
	}

	// ")"
	if p.index >= len(p.tokens) && p.tokens[p.index].Kind != tokenizer.SRPAREN {
		return errors.New(fmt.Sprintf("syntax error: cannot find ')'"))
	}

	p.index++

	// "{"
	if p.index >= len(p.tokens) && p.tokens[p.index].Kind != tokenizer.SLBRACE {
		return errors.New(fmt.Sprintf("syntax error: cannot find '{'"))
	}
	
	p.index++

	// 記述部
	err = p.description()
	if err != nil {
		return err
	}

	// "}"
	if p.index >= len(p.tokens) && p.tokens[p.index].Kind != tokenizer.SRBRACE {
		return errors.New(fmt.Sprintf("syntax error: cannot find '}'"))
	}

	p.index++
	
	return nil
}

// else節
func (p *Parser) elseSection() error {
	var err error

	// "else"
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != tokenizer.SELSE {
		return errors.New(fmt.Sprintf("syntax error: cannot find 'else'"))
	}

	p.index++

	// "{"
	if p.index >= len(p.tokens) && p.tokens[p.index].Kind != tokenizer.SLBRACE {
		return errors.New(fmt.Sprintf("syntax error: cannot find '{'"))
	}
	
	p.index++

	// 記述部
	err = p.description()
	if err != nil {
		return err
	}

	// "}"
	if p.index >= len(p.tokens) && p.tokens[p.index].Kind != tokenizer.SRBRACE {
		return errors.New(fmt.Sprintf("syntax error: cannot find '}'"))
	}

	p.index++
	
	return nil
}

// 関数名
func (p *Parser) functionName() error {
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != tokenizer.SIDENTIFIER {
		return errors.New(fmt.Sprintf("syntax error: cannot find an identifier"))
	}

	p.index++

	return nil
}

// 変数名
func (p *Parser) variableName() error {
	if p.index >= len(p.tokens) || p.tokens[p.index].Kind != tokenizer.SIDENTIFIER {
		return errors.New(fmt.Sprintf("syntax error: cannot find an identifier"))
	}

	p.index++

	return nil
}
