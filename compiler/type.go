package compiler

import (
	"fmt"
	"errors"

	"myriad/tokenizer"
)

// 中間言語
type IntefCode interface {
	getName() string
}

type InterCodeKind int
const (
	ROW InterCodeKind = iota
	COMMAND
	VAR
	IF
	ELIF
	ELSE
	ENDIF
)

type InterCodeCommonDetail struct {
	Content string
	Kind InterCodeKind
}

// ROW, COMMAND, VAR, ELSE, ENDIF
type NormalInterCode struct {
	InterCodeCommonDetail
}

type OperaterKind int
const (
	EQUAL OperaterKind = iota
	NOTEQUAL
)

type Formula struct {
	Content string
	Kind tokenizer.TokenKind
}

type IfContent struct {
	LFormula, RFormula Formula
	Operator OperaterKind
	NextOffset int
	EndOffset int
}

// IF, ELIF
type IfInterCode struct {
	InterCodeCommonDetail
	IfContent
}

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

// 単一変数用のインタフェース実装
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

// 配列型変数用のインタフェース実装
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
	intToString := map[InterCodeKind]string{
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