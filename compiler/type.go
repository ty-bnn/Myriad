package compiler

import (
	"errors"
	"fmt"

	"myriad/tokenizer"
)

// InterCode 中間言語
type InterCode interface {
	GetContent() string
	GetKind() int
}

const (
	ROW int = iota
	COMMAND
	VAR
	IF
	ELIF
	ELSE
	ENDIF
)

// InterCodeCommonDetail 中間言語の共通部分
type InterCodeCommonDetail struct {
	Content string
	Kind    int
}

// NormalInterCode ROW, COMMAND, VAR, ELSE, ENDIF
type NormalInterCode struct {
	InterCodeCommonDetail
}

func (n NormalInterCode) GetContent() string {
	return n.Content
}

func (n NormalInterCode) GetKind() int {
	return n.Kind
}

const (
	EQUAL int = iota
	NOTEQUAL
)

type Formula struct {
	Content string
	Kind    tokenizer.TokenKind
}

type IfContent struct {
	LFormula, RFormula Formula
	Operator           int
	NextOffset         int
	EndOffset          int
}

// IfInterCode IF, ELIF
type IfInterCode struct {
	InterCodeCommonDetail
	IfContent
}

func (i IfInterCode) GetContent() string {
	return i.Content
}

func (i IfInterCode) GetKind() int {
	return i.Kind
}

// Variable TODO: データ構造再考
/*
(注) GetKind, getName
Variable型でSingleVariableやMultiVariableを扱う場合、フィールド変数にアクセスできない
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

// SingleVariable 単一変数用のインタフェース実装
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

// MultiVariable 配列型変数用のインタフェース実装
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
