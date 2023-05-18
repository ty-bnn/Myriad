package parser

import (
	"fmt"
	"errors"

	"dcc/tokenizer"
)

func (p *Parser) program(index int) error {
	var err error

	// { 関数インポート文 }
	for ;; {
		if index >= len(p.tokens) || p.tokens[index].Kind != tokenizer.SIMPORT {
			break
		}

		index, err = p.importFunc(index)
		if err != nil {
			return err
		}
	}

	// { 関数 }
	for ;; {
		if index >= len(p.tokens) || p.tokens[index].Kind != tokenizer.SIDENTIFIER {
			break
		}

		index, err = p.function(index)
		if err != nil {
			return err
		}
	}

	if index >= len(p.tokens) || p.tokens[index].Kind != tokenizer.SMAIN {
		return nil
	}
	
	// メイン部
	index, err = p.mainFunction(index)
	if err != nil {
		return err
	}

	return nil
}

// 関数インポート文
func (p *Parser) importFunc(index int) (int, error) {
	var err error
	// "import"
	if index >= len(p.tokens) || p.tokens[index].Kind != tokenizer.SIMPORT {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find 'import'"))
	}

	index++

	// 関数名
	index, err = p.functionName(index)
	if err != nil {
		return index, err
	}

	// "from"
	if index >= len(p.tokens) || p.tokens[index].Kind != tokenizer.SFROM {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find 'from'"))
	}

	index++

	// ファイル名
	index, err = p.fileName(index)
	if err != nil {
		return index, err
	}

	return index, nil
}

// ファイル名
func (p *Parser) fileName(index int) (int, error) {
	if index >= len(p.tokens) || p.tokens[index].Kind != tokenizer.SSTRING {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find an identifier"))
	}

	index++

	return index, nil
}

// 関数
func (p *Parser) function(index int) (int, error) {
	var err error

	// 関数名
	index, err = p.functionName(index)
	if err != nil {
		return index, err
	}

	// 引数宣言部
	index, err = p.argumentDecralation(index)
	if err != nil {
		return index, err
	}

	// 関数記述部
	index, err = p.functionDescription(index)
	if err != nil {
		return index, err
	}

	return index, nil
}

// メイン部
func (p *Parser) mainFunction(index int) (int, error) {
	var err error

	// "main"
	if index >= len(p.tokens) || p.tokens[index].Kind != tokenizer.SMAIN {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find 'main' to declare"))
	}

	index++

	// 引数宣言部
	index, err = p.argumentDecralation(index)
	if err != nil {
		return index, err
	}

	// 関数記述部
	index, err = p.functionDescription(index)
	if err != nil {
		return index, err
	}

	return index, nil
}

// 引数宣言部
func (p *Parser) argumentDecralation(index int) (int, error) {
	var err error

	// "("
	if index >= len(p.tokens) || p.tokens[index].Kind != tokenizer.SLPAREN {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find '('"))
	}

	index++

	// 引数群
	if 	index < len(p.tokens) && p.tokens[index].Kind == tokenizer.SIDENTIFIER {
		index, err = p.arguments(index)
		if err != nil {
			return index, err
		}
	}

	// ")"
	if index >= len(p.tokens) || p.tokens[index].Kind != tokenizer.SRPAREN {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find ')'"))
	}

	index++

	return index, nil
}

// 引数群
func (p *Parser) arguments(index int) (int, error) {
	var err error

	// 変数
	index, err = p.variable(index)
	if err != nil {
		return index, err
	}

	for ;; {
		// ","
		if index >= len(p.tokens) || p.tokens[index].Kind != tokenizer.SCOMMA {
			break
		}

		index++
		
		// 変数
		index, err = p.variable(index)
		if err != nil {
			return index, err
		}
	}

	return index, nil
}

// 変数
func (p *Parser) variable(index int) (int, error) {
	var err error

	// 変数名
	index, err = p.variableName(index)
	if err != nil {
		return index, err
	}

	// "[]"
	if index < len(p.tokens) && p.tokens[index].Kind == tokenizer.SARRANGE {
		index++
	}

	return index, nil
}

// 関数記述部
func (p *Parser) functionDescription(index int) (int, error) {
	var err error

	// "{"
	if index >= len(p.tokens) || p.tokens[index].Kind != tokenizer.SLBRACE {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find '{'"))
	}

	index++

	// 記述部
	index, err = p.description(index)
	if err != nil {
		return index, err
	}

	// "}"
	if index >= len(p.tokens) || p.tokens[index].Kind != tokenizer.SRBRACE {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find '}'"))
	}

	index++

	return index, nil
}

// 記述部
func (p *Parser) description(index int) (int, error) {
	var err error

	// 記述ブロック
	index, err = p.descriptionBlock(index)
	if err != nil {
		return index, err
	}

	for ;; {
		if index >= len(p.tokens) || (p.tokens[index].Kind != tokenizer.SDFCOMMAND && p.tokens[index].Kind != tokenizer.SDFARG && p.tokens[index].Kind != tokenizer.SIDENTIFIER && p.tokens[index].Kind != tokenizer.SIF) {
			break
		}
			
		index, err = p.descriptionBlock(index)
		if err != nil {
			return index, err
		}
	}

	return index, nil
}

// 記述ブロック
func (p *Parser) descriptionBlock(index int) (int, error) {
	var err error

	if index < len(p.tokens) && p.tokens[index].Kind == tokenizer.SDFCOMMAND || p.tokens[index].Kind == tokenizer.SDFARG {
		// Dfile文
		index, err = p.dockerFile(index)
		if err != nil {
			return index, err
		}
	} else if index + 1 < len(p.tokens) && p.tokens[index].Kind == tokenizer.SIDENTIFIER && p.tokens[index + 1].Kind == tokenizer.SLPAREN {
		// 関数呼び出し文
		index, err = p.functionCall(index)
		if err != nil {
			return index, err
		}
	} else if index + 1 < len(p.tokens) && p.tokens[index].Kind == tokenizer.SIDENTIFIER && p.tokens[index + 1].Kind == tokenizer.SDEFINE {
		// 変数定義文
		index, err = p.defineVariable(index)
	} else if index < len(p.tokens) && p.tokens[index].Kind == tokenizer.SIF {
		// ifブロック
		index, err = p.ifBlock(index)
		if err != nil {
			return index, err
		}
	} else {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find a description block"))
	}

	return index, nil
}

// Dfile文
func (p *Parser) dockerFile(index int) (int, error) {
	var err error
	// Df命令
	if index >= len(p.tokens) || (p.tokens[index].Kind != tokenizer.SDFCOMMAND && p.tokens[index].Kind != tokenizer.SDFARG) {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find a Dockerfile comamnd"))
	}

	if p.tokens[index].Kind == tokenizer.SDFCOMMAND {
		index++
	}


	// Df引数部
	index, err = p.dfArgs(index)
	if err != nil {
		return index, err
	}

	return index, nil
}

// Df引数部
func (p *Parser) dfArgs(index int) (int, error) {
	var err error
	index, err = p.dfArg(index)
	if err != nil {
		return index, err
	}

	for ;; {
		if index < len(p.tokens) && (p.tokens[index].Kind == tokenizer.SDFARG || p.tokens[index].Kind == tokenizer.SASSIGNVARIABLE) {
			index, err = p.dfArg(index)
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
func (p *Parser) dfArg(index int) (int, error) {
	if index >= len(p.tokens) || (p.tokens[index].Kind != tokenizer.SDFARG && p.tokens[index].Kind != tokenizer.SASSIGNVARIABLE) {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find Df argument"))
	}

	index++

	return index, nil
}

// 関数呼び出し文
func (p *Parser) functionCall(index int) (int, error) {
	var err error

	// 関数名
	index, err = p.functionName(index)
	if err != nil {
		return index, err
	}

	// "("
	if index >= len(p.tokens) || p.tokens[index].Kind != tokenizer.SLPAREN {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find '('"))
	}

	index++

	// 式の並び
	if index < len(p.tokens) && (p.tokens[index].Kind == tokenizer.SSTRING || p.tokens[index].Kind == tokenizer.SIDENTIFIER) {
		index, err = p.rowOfFormulas(index)
		if err != nil {
			return index, err
		}
	}

	// ")"
	if index >= len(p.tokens) || p.tokens[index].Kind != tokenizer.SRPAREN {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find ')'"))
	}

	index++

	return index, nil
}

// 式の並び
func (p *Parser) rowOfFormulas(index int) (int, error) {
	var err error
	// 式
	index, err = p.formula(index)
	if err != nil {
		return index, err
	}

	for ;; {
		// ","
		if index >= len(p.tokens) || p.tokens[index].Kind != tokenizer.SCOMMA {
			break
		}

		index++


		// 式
		index, err = p.formula(index)
		if err != nil {
			return index, err
		}
	}

	return index, nil
}

// ifブロック
func (p *Parser) ifBlock(index int) (int, error) {
	var err error

	// "if"
	if index >= len(p.tokens) && p.tokens[index].Kind != tokenizer.SIF {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find 'if'"))
	}

	index++

	// "("
	if index >= len(p.tokens) && p.tokens[index].Kind != tokenizer.SLPAREN {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find '('"))
	}

	index++

	// 条件判定式
	index, err = p.conditionalFormula(index)
	if err != nil {
		return index, err
	}

	// ")"
	if index >= len(p.tokens) && p.tokens[index].Kind != tokenizer.SRPAREN {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find ')'"))
	}

	index++

	// "{"
	if index >= len(p.tokens) && p.tokens[index].Kind != tokenizer.SLBRACE {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find '{'"))
	}
	
	index++

	// 記述部
	index, err = p.description(index)
	if err != nil {
		return index, err
	}

	// "}"
	if index >= len(p.tokens) && p.tokens[index].Kind != tokenizer.SRBRACE {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find '}'"))
	}

	index++

	for ;; {
		if index >= len(p.tokens) || p.tokens[index].Kind != tokenizer.SELIF {
			break
		}

		// elif節
		index, err = p.elifSection(index)
		if err != nil {
			return index, err
		}
	}

	if index < len(p.tokens) && p.tokens[index].Kind == tokenizer.SELSE {
		index, err = p.elseSection(index)
		if err != nil {
			return index, err
		}
	}

	return index, nil
}

// 条件判定式
func (p *Parser) conditionalFormula(index int) (int, error) {
	var err error
	// 式
	index, err = p.formula(index)
	if err != nil {
		return index, err
	}

	// 比較演算子
	index, err = p.conditionalOperator(index)
	if err != nil {
		return index, err
	}

	// 式
	index, err = p.formula(index)
	if err != nil {
		return index, err
	}

	return index, nil
}

// 式
func (p *Parser) formula(index int) (int, error) {
	// 変数, 文字列
	if index >= len(p.tokens) || (p.tokens[index].Kind != tokenizer.SIDENTIFIER && p.tokens[index].Kind != tokenizer.SSTRING) {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find a formula"))
	}

	index++

	return index, nil
}

// 比較演算子
func (p *Parser) conditionalOperator(index int) (int, error) {
	// "==", "!="
	if index >= len(p.tokens) || (p.tokens[index].Kind != tokenizer.SEQUAL && p.tokens[index].Kind != tokenizer.SNOTEQUAL) {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find conditional operator"))
	}

	index++

	return index, nil
}

// 変数定義文
func (p *Parser) defineVariable(index int) (int, error) {
	var err error

	// 変数名
	index, err = p.variableName(index)
	if err != nil {
		return index, err
	}

	// ":="
	if index >= len(p.tokens) || p.tokens[index].Kind != tokenizer.SDEFINE {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find ':='"))
	}

	index++

	// 文字列
	if index >= len(p.tokens) || p.tokens[index].Kind != tokenizer.SSTRING {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find string"))
	}

	index++

	return index, nil
}

// elif節
func (p *Parser) elifSection(index int) (int, error) {
	var err error

	// "else if"
	if index >= len(p.tokens) || p.tokens[index].Kind != tokenizer.SELIF {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find 'else if'"))
	}

	index++

	// "("
	if index >= len(p.tokens) && p.tokens[index].Kind != tokenizer.SLPAREN {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find '('"))
	}

	index++

	// 条件判定式
	index, err = p.conditionalFormula(index)
	if err != nil {
		return index, err
	}

	// ")"
	if index >= len(p.tokens) && p.tokens[index].Kind != tokenizer.SRPAREN {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find ')'"))
	}

	index++

	// "{"
	if index >= len(p.tokens) && p.tokens[index].Kind != tokenizer.SLBRACE {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find '{'"))
	}
	
	index++

	// 記述部
	index, err = p.description(index)
	if err != nil {
		return index, err
	}

	// "}"
	if index >= len(p.tokens) && p.tokens[index].Kind != tokenizer.SRBRACE {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find '}'"))
	}

	index++
	
	return index, nil
}

// else節
func (p *Parser) elseSection(index int) (int, error) {
	var err error

	// "else"
	if index >= len(p.tokens) || p.tokens[index].Kind != tokenizer.SELSE {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find 'else'"))
	}

	index++

	// "{"
	if index >= len(p.tokens) && p.tokens[index].Kind != tokenizer.SLBRACE {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find '{'"))
	}
	
	index++

	// 記述部
	index, err = p.description(index)
	if err != nil {
		return index, err
	}

	// "}"
	if index >= len(p.tokens) && p.tokens[index].Kind != tokenizer.SRBRACE {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find '}'"))
	}

	index++
	
	return index, nil
}

// 関数名
func (p *Parser) functionName(index int) (int, error) {
	if index >= len(p.tokens) || p.tokens[index].Kind != tokenizer.SIDENTIFIER {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find an identifier"))
	}

	index++

	return index, nil
}

// 変数名
func (p *Parser) variableName(index int) (int, error) {
	if index >= len(p.tokens) || p.tokens[index].Kind != tokenizer.SIDENTIFIER {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find an identifier"))
	}

	index++

	return index, nil
}
