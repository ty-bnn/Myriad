package compiler

import (
	"fmt"
	"errors"

	"myriad/tokenizer"
)

// 中間言語
type InterCode struct {
	Content string
	Kind CodeKind
	ArgValues []string // 関数呼び出し
	IfContent IfContent // if文
}

type CodeKind int
const (
	ROW CodeKind = iota
	COMMAND
	VAR
	IF
	ELIF
	ELSE
	ENDIF
)

type IfContent struct {
	LFormula, RFormula Formula
	Operator OperaterKind
	NextOffset int
	EndOffset int
}

type Formula struct {
	Content string
	Kind tokenizer.TokenKind
}

type OperaterKind int
const (
	EQUAL OperaterKind = iota
	NOTEQUAL
)

// TODO: データ構造再考
/*
(注) getKind, getName
Variable型で変数を扱う場合、フィールド変数にアクセスできない
SingleVariableとMultiVariableを共通してVariable型で扱う場合のみ使用する。
それ以外は直接フィールド変数にアクセス可能
*/
type Variable interface {
	getValue(int) (string, error)
	getName() string
	getKind() VariableKind
}

type VariableKind int
const (
	VARIABLE VariableKind = iota
	ARGUMENT
)

type VariableCommonDetail struct {
	Name string
	Kind VariableKind
}

type SingleVariable struct {
	VariableCommonDetail
	Value string
}

func (s SingleVariable) getValue(index int) (string, error) {
	if index != 0 {
		return "", errors.New(fmt.Sprintf("semantic error: single values do not have index"))
	}

	return s.Value, nil
}

func (s SingleVariable) getName() string {
	return s.Name
}

func (s SingleVariable) getKind() VariableKind {
	return s.Kind
}

type MultiVariable struct {
	VariableCommonDetail
	Values []string
}

func (m MultiVariable) getValue(index int) (string, error) {
	if index < 0 || len(m.Values) <= index {
		return "", errors.New(fmt.Sprintf("semantic error: out of index"))
	}

	return m.Values[index], nil
}

func (m MultiVariable) getName() string {
	return m.Name
}

func (m MultiVariable) getKind() VariableKind {
	return m.Kind
}

func (c Compiler) printInterCodes(functionName string) {
	intToString := map[CodeKind]string{
		ROW: "row",
		COMMAND: "command",
		VAR: "var",
		IF: "if",
		ELIF: "elif",
		ELSE: "else",
		ENDIF: "endif",
	}

	fmt.Printf("--------- inter codes in \"%s\"---------\n", functionName)
	for i, code := range c.FunctionInterCodeMap[functionName] {
		fmt.Println("{")
		fmt.Printf("  Content: %s\n", code.Content)
		fmt.Printf("     Kind: %s\n", intToString[code.Kind])
		fmt.Printf("    Index: %2d\n", i)
		fmt.Println("},")
	}
}