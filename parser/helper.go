package parser

import (
	"fmt"
	"errors"

	"dcc/tokenizer"
)

func program(tokens []tokenizer.Token, index int) error {
	var err error

	// { 関数インポート文 }
	for ;; {
		if index >= len(tokens) || tokens[index].Kind != tokenizer.SIMPORT {
			break
		}

		index, err = importFunc(tokens, index)
		if err != nil {
			return err
		}
	}

	// { 関数 }
	for ;; {
		if index >= len(tokens) || tokens[index].Kind != tokenizer.SIDENTIFIER {
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
	fmt.Println("hellllo")
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
	if index >= len(tokens) || tokens[index].Kind != tokenizer.SIMPORT {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find 'import'"))
	}

	index++

	// 関数名
	index, err = functionName(tokens, index)
	if err != nil {
		return index, err
	}

	// "from"
	if index >= len(tokens) || tokens[index].Kind != tokenizer.SFROM {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find 'from'"))
	}

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
	if index >= len(tokens) || tokens[index].Kind != tokenizer.SSTRING {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find an identifier"))
	}

	index++

	return index, nil
}

// 関数
func function(tokens []tokenizer.Token, index int) (int, error) {
	var err error

	// 関数名
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

	// "main"
	if index >= len(tokens) || tokens[index].Kind != tokenizer.SMAIN {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find 'main' to declare"))
	}

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
	if index >= len(tokens) || tokens[index].Kind != tokenizer.SLPAREN {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find '('"))
	}

	index++

	// 引数群
	if 	index < len(tokens) && tokens[index].Kind == tokenizer.SIDENTIFIER {
		index, err = arguments(tokens, index)
		if err != nil {
			return index, err
		}
	}

	// ")"
	if index >= len(tokens) || tokens[index].Kind != tokenizer.SRPAREN {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find ')'"))
	}

	index++

	return index, nil
}

// 引数群
func arguments(tokens []tokenizer.Token, index int) (int, error) {
	var err error

	// 変数
	index, err = variable(tokens, index)
	if err != nil {
		return index, err
	}

	for ;; {
		// ","
		if index >= len(tokens) || tokens[index].Kind != tokenizer.SCOMMA {
			break
		}

		index++
		
		// 変数
		index, err = variable(tokens, index)
		if err != nil {
			return index, err
		}
	}

	return index, nil
}

// 変数
func variable(tokens []tokenizer.Token, index int) (int, error) {
	var err error

	// 変数名
	index, err = variableName(tokens, index)
	if err != nil {
		return index, err
	}

	// "[]"
	if index < len(tokens) && tokens[index].Kind == tokenizer.SARRANGE {
		index++
	}

	return index, nil
}

// 関数記述部
func functionDescription(tokens []tokenizer.Token, index int) (int, error) {
	var err error

	// "{"
	if index >= len(tokens) || tokens[index].Kind != tokenizer.SLBRACE {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find '{'"))
	}

	index++

	// 記述部
	index, err = description(tokens, index)
	if err != nil {
		return index, err
	}

	// "}"
	if index >= len(tokens) || tokens[index].Kind != tokenizer.SRBRACE {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find '}'"))
	}

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
		if index >= len(tokens) || (tokens[index].Kind != tokenizer.SDFCOMMAND && tokens[index].Kind != tokenizer.SDFARG && tokens[index].Kind != tokenizer.SIDENTIFIER && tokens[index].Kind != tokenizer.SIF) {
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

	if index < len(tokens) && tokens[index].Kind == tokenizer.SDFCOMMAND || tokens[index].Kind == tokenizer.SDFARG {
		// Dfile文
		index, err = dockerFile(tokens, index)
		if err != nil {
			return index, err
		}
	} else if index + 1 < len(tokens) && tokens[index].Kind == tokenizer.SIDENTIFIER && tokens[index + 1].Kind == tokenizer.SLPAREN {
		// 関数呼び出し文
		index, err = functionCall(tokens, index)
		if err != nil {
			return index, err
		}
	} else if index + 1 < len(tokens) && tokens[index].Kind == tokenizer.SIDENTIFIER && tokens[index + 1].Kind == tokenizer.SDEFINE {
		// 変数定義文
		index, err = defineVariable(tokens, index)
	} else if index < len(tokens) && tokens[index].Kind == tokenizer.SIF {
		// ifブロック
		index, err = ifBlock(tokens, index)
		if err != nil {
			return index, err
		}
	} else {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find a description block"))
	}

	return index, nil
}

// Dfile文
func dockerFile(tokens []tokenizer.Token, index int) (int, error) {
	var err error
	// Df命令
	if index >= len(tokens) || (tokens[index].Kind != tokenizer.SDFCOMMAND && tokens[index].Kind != tokenizer.SDFARG) {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find a Dockerfile comamnd"))
	}

	if tokens[index].Kind == tokenizer.SDFCOMMAND {
		index++
	}


	// Df引数部
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
		if index < len(tokens) && (tokens[index].Kind == tokenizer.SDFARG || tokens[index].Kind == tokenizer.SASSIGNVARIABLE) {
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
	if index >= len(tokens) || (tokens[index].Kind != tokenizer.SDFARG && tokens[index].Kind != tokenizer.SASSIGNVARIABLE) {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find Df argument"))
	}

	index++

	return index, nil
}

// 関数呼び出し文
func functionCall(tokens []tokenizer.Token, index int) (int, error) {
	var err error

	// 関数名
	index, err = functionName(tokens, index)
	if err != nil {
		return index, err
	}

	// "("
	if index >= len(tokens) || tokens[index].Kind != tokenizer.SLPAREN {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find '('"))
	}

	index++

	// 式の並び
	if index < len(tokens) && (tokens[index].Kind == tokenizer.SSTRING || tokens[index].Kind == tokenizer.SIDENTIFIER) {
		index, err = rowOfFormulas(tokens, index)
		if err != nil {
			return index, err
		}
	}

	// ")"
	if index >= len(tokens) || tokens[index].Kind != tokenizer.SRPAREN {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find ')'"))
	}

	index++

	return index, nil
}

// 式の並び
func rowOfFormulas(tokens []tokenizer.Token, index int) (int, error) {
	var err error
	// 式
	index, err = formula(tokens, index)
	if err != nil {
		return index, err
	}

	for ;; {
		// ","
		if index >= len(tokens) || tokens[index].Kind != tokenizer.SCOMMA {
			break
		}

		index++


		// 式
		index, err = formula(tokens, index)
		if err != nil {
			return index, err
		}
	}

	return index, nil
}

// ifブロック
func ifBlock(tokens []tokenizer.Token, index int) (int, error) {
	var err error

	// "if"
	if index >= len(tokens) && tokens[index].Kind != tokenizer.SIF {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find 'if'"))
	}

	index++

	// "("
	if index >= len(tokens) && tokens[index].Kind != tokenizer.SLPAREN {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find '('"))
	}

	index++

	// 条件判定式
	index, err = conditionalFormula(tokens, index)
	if err != nil {
		return index, err
	}

	// ")"
	if index >= len(tokens) && tokens[index].Kind != tokenizer.SRPAREN {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find ')'"))
	}

	index++

	// "{"
	if index >= len(tokens) && tokens[index].Kind != tokenizer.SLBRACE {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find '{'"))
	}
	
	index++

	// 記述部
	index, err = description(tokens, index)
	if err != nil {
		return index, err
	}

	// "}"
	if index >= len(tokens) && tokens[index].Kind != tokenizer.SRBRACE {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find '}'"))
	}

	index++

	for ;; {
		if index >= len(tokens) || tokens[index].Kind != tokenizer.SELIF {
			break
		}

		// elif節
		index, err = elifSection(tokens, index)
		if err != nil {
			return index, err
		}
	}

	if index < len(tokens) && tokens[index].Kind == tokenizer.SELSE {
		index, err = elseSection(tokens, index)
		if err != nil {
			return index, err
		}
	}

	return index, nil
}

// 条件判定式
func conditionalFormula(tokens []tokenizer.Token, index int) (int, error) {
	var err error
	// 式
	index, err = formula(tokens, index)
	if err != nil {
		return index, err
	}

	// 比較演算子
	index, err = conditionalOperator(tokens, index)
	if err != nil {
		return index, err
	}

	// 式
	index, err = formula(tokens, index)
	if err != nil {
		return index, err
	}

	return index, nil
}

// 式
func formula(tokens []tokenizer.Token, index int) (int, error) {
	// 変数, 文字列
	if index >= len(tokens) || (tokens[index].Kind != tokenizer.SIDENTIFIER && tokens[index].Kind != tokenizer.SSTRING) {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find a formula"))
	}

	index++

	return index, nil
}

// 比較演算子
func conditionalOperator(tokens []tokenizer.Token, index int) (int, error) {
	// "==", "!="
	if index >= len(tokens) || (tokens[index].Kind != tokenizer.SEQUAL && tokens[index].Kind != tokenizer.SNOTEQUAL) {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find conditional operator"))
	}

	index++

	return index, nil
}

// 変数定義文
func defineVariable(tokens []tokenizer.Token, index int) (int, error) {
	var err error

	// 変数名
	index, err = variableName(tokens, index)
	if err != nil {
		return index, err
	}

	// ":="
	if index >= len(tokens) || tokens[index].Kind != tokenizer.SDEFINE {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find ':='"))
	}

	index++

	// 文字列
	if index >= len(tokens) || tokens[index].Kind != tokenizer.SSTRING {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find string"))
	}

	index++

	return index, nil
}

// elif節
func elifSection(tokens []tokenizer.Token, index int) (int, error) {
	var err error

	// "else if"
	if index >= len(tokens) || tokens[index].Kind != tokenizer.SELIF {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find 'else if'"))
	}

	index++

	// "("
	if index >= len(tokens) && tokens[index].Kind != tokenizer.SLPAREN {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find '('"))
	}

	index++

	// 条件判定式
	index, err = conditionalFormula(tokens, index)
	if err != nil {
		return index, err
	}

	// ")"
	if index >= len(tokens) && tokens[index].Kind != tokenizer.SRPAREN {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find ')'"))
	}

	index++

	// "{"
	if index >= len(tokens) && tokens[index].Kind != tokenizer.SLBRACE {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find '{'"))
	}
	
	index++

	// 記述部
	index, err = description(tokens, index)
	if err != nil {
		return index, err
	}

	// "}"
	if index >= len(tokens) && tokens[index].Kind != tokenizer.SRBRACE {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find '}'"))
	}

	index++
	
	return index, nil
}

// else節
func elseSection(tokens []tokenizer.Token, index int) (int, error) {
	var err error

	// "else"
	if index >= len(tokens) || tokens[index].Kind != tokenizer.SELSE {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find 'else'"))
	}

	index++

	// "{"
	if index >= len(tokens) && tokens[index].Kind != tokenizer.SLBRACE {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find '{'"))
	}
	
	index++

	// 記述部
	index, err = description(tokens, index)
	if err != nil {
		return index, err
	}

	// "}"
	if index >= len(tokens) && tokens[index].Kind != tokenizer.SRBRACE {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find '}'"))
	}

	index++
	
	return index, nil
}

// 関数名
func functionName(tokens []tokenizer.Token, index int) (int, error) {
	if index >= len(tokens) || tokens[index].Kind != tokenizer.SIDENTIFIER {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find an identifier"))
	}

	index++

	return index, nil
}

// 変数名
func variableName(tokens []tokenizer.Token, index int) (int, error) {
	if index >= len(tokens) || tokens[index].Kind != tokenizer.SIDENTIFIER {
		return index, errors.New(fmt.Sprintf("syntax error: cannot find an identifier"))
	}

	index++

	return index, nil
}
